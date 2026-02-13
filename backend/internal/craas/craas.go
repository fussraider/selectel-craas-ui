package craas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"sync"
	"time"

	v1 "github.com/selectel/craas-go/pkg"
	"github.com/selectel/craas-go/pkg/v1/registry"
	"github.com/selectel/craas-go/pkg/v1/repository"
)

const Endpoint = "https://cr.selcloud.ru/api/v1"

type Service struct {
	endpoint string
	logger   *slog.Logger
}

func New(logger *slog.Logger) *Service {
	return &Service{
		endpoint: Endpoint,
		logger:   logger.With("service", "craas"),
	}
}

// ListRegistries returns a list of registries for the project (scoped by token).
func (s *Service) ListRegistries(ctx context.Context, token string) ([]*registry.Registry, error) {
	s.logger.Debug("listing registries")
	client := v1.NewCRaaSClientV1(token, s.endpoint)

	start := time.Now()
	registries, _, err := registry.List(ctx, client)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to list registries", "error", err)
		return nil, err
	}
	s.logger.Info("listed registries", "count", len(registries), "duration", duration)
	return registries, nil
}

// ListRepositories returns a list of repositories in the registry.
func (s *Service) ListRepositories(ctx context.Context, token string, registryID string) ([]*repository.Repository, error) {
	s.logger.Debug("listing repositories", "registry_id", registryID)
	client := v1.NewCRaaSClientV1(token, s.endpoint)

	start := time.Now()
	repos, _, err := repository.ListRepositories(ctx, client, registryID)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to list repositories", "registry_id", registryID, "error", err)
		return nil, err
	}
	s.logger.Info("listed repositories", "registry_id", registryID, "count", len(repos), "duration", duration)
	return repos, nil
}

// ListImages returns a list of images in the repository.
func (s *Service) ListImages(ctx context.Context, token string, registryID, repoName string) ([]*repository.Image, error) {
	s.logger.Debug("listing images", "registry_id", registryID, "repository", repoName)
	client := v1.NewCRaaSClientV1(token, s.endpoint)

	start := time.Now()
	images, _, err := repository.ListImages(ctx, client, registryID, repoName)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to list images", "registry_id", registryID, "repository", repoName, "error", err)
		return nil, err
	}

	// Fetch all tags to check for missing ones
	allTags, _, err := repository.ListTags(ctx, client, registryID, repoName)
	if err != nil {
		s.logger.Warn("failed to list tags for verification", "registry_id", registryID, "repository", repoName, "error", err)
		return images, nil
	}

	// Create a map of existing images by digest
	imageMap := make(map[string]*repository.Image)
	knownTags := make(map[string]struct{})
	for _, img := range images {
		imageMap[img.Digest] = img
		for _, t := range img.Tags {
			knownTags[t] = struct{}{}
		}
	}

	// Identify missing tags
	var missingTags []string
	for _, t := range allTags {
		if _, ok := knownTags[t]; !ok {
			missingTags = append(missingTags, t)
		}
	}

	if len(missingTags) > 0 {
		s.logger.Info("found missing tags, resolving", "count", len(missingTags), "tags", missingTags)

		// Resolve missing tags concurrently
		var wg sync.WaitGroup
		sem := make(chan struct{}, 5) // Limit to 5 concurrent requests
		var mu sync.Mutex

		httpClient := &http.Client{
			Timeout: 10 * time.Second,
		}

		for _, tag := range missingTags {
			wg.Add(1)
			sem <- struct{}{} // Acquire token
			go func(tag string) {
				defer wg.Done()
				defer func() { <-sem }() // Release token

				// Fetch all digests associated with the tag
				digests, err := s.fetchImageDigests(ctx, httpClient, token, registryID, repoName, tag)
				if err != nil {
					s.logger.Warn("failed to fetch digests for tag", "tag", tag, "error", err)
					return
				}

				mu.Lock()
				defer mu.Unlock()

				for _, digest := range digests {
					// Check if we already have this image (by digest)
					if existingImg, ok := imageMap[digest]; ok {
						// Add the tag to the existing image
						exists := false
						for _, t := range existingImg.Tags {
							if t == tag {
								exists = true
								break
							}
						}
						if !exists {
							existingImg.Tags = append(existingImg.Tags, tag)
						}
					} else {
						// New image found with this digest.
						// We create a minimal entry for it.
						newImg := &repository.Image{
							Digest:    digest,
							Tags:      []string{tag},
							CreatedAt: time.Now(), // Unknown
						}
						images = append(images, newImg)
						imageMap[digest] = newImg
					}
				}
			}(tag)
		}
		wg.Wait()
	}

	s.logger.Info("listed images", "registry_id", registryID, "repository", repoName, "count", len(images), "duration", duration)
	return images, nil
}

// digestRegex matches standard SHA256 digests.
var digestRegex = regexp.MustCompile(`^sha256:[a-f0-9]{64}$`)

// fetchImageDigests fetches the digest(s) associated with a tag.
// Returns a slice of digests found in the response (header or body).
// It recursively searches the JSON body for any string that matches a SHA256 digest format.
func (s *Service) fetchImageDigests(ctx context.Context, client *http.Client, token, registryID, repoName, reference string) ([]string, error) {
	url := fmt.Sprintf("%s/registries/%s/repositories/%s/%s", s.endpoint, registryID, repoName, reference)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", token)
	// Add Accept headers to request Manifests/Indices properly instead of empty layer lists
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json, application/vnd.docker.distribution.manifest.list.v2+json, application/vnd.oci.image.manifest.v1+json, application/vnd.oci.image.index.v1+json, */*")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var digests []string
	digestSet := make(map[string]struct{})

	// 1. Check Header
	headerDigest := resp.Header.Get("Docker-Content-Digest")
	if headerDigest != "" {
		digests = append(digests, headerDigest)
		digestSet[headerDigest] = struct{}{}
	}

	// 2. Deep search in JSON body
	if len(bodyBytes) > 0 {
		var data interface{}
		if err := json.Unmarshal(bodyBytes, &data); err == nil {
			findDigests(data, digestSet, &digests)
		} else {
			// If not valid JSON, we can't search it structurally.
			// Maybe regex search on raw bytes if desperate?
			// Let's stick to structural search for safety first.
		}
	}

	if len(digests) == 0 {
		// Log truncated body for debugging
		maxLog := 512
		if len(bodyBytes) < maxLog {
			maxLog = len(bodyBytes)
		}
		s.logger.Warn("no digests found in response body", "tag", reference, "body_preview", string(bodyBytes[:maxLog]))
		return nil, fmt.Errorf("no digests found for tag %s", reference)
	}

	return digests, nil
}

// findDigests recursively searches for strings matching digestRegex in the JSON structure.
func findDigests(data interface{}, seen map[string]struct{}, result *[]string) {
	switch v := data.(type) {
	case string:
		if digestRegex.MatchString(v) {
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				*result = append(*result, v)
			}
		}
	case []interface{}:
		for _, item := range v {
			findDigests(item, seen, result)
		}
	case map[string]interface{}:
		for _, value := range v {
			findDigests(value, seen, result)
		}
	}
}

// ListTags returns a list of tags in the repository.
func (s *Service) ListTags(ctx context.Context, token string, registryID, repoName string) ([]string, error) {
	s.logger.Debug("listing tags", "registry_id", registryID, "repository", repoName)
	client := v1.NewCRaaSClientV1(token, s.endpoint)

	start := time.Now()
	tags, _, err := repository.ListTags(ctx, client, registryID, repoName)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to list tags", "registry_id", registryID, "repository", repoName, "error", err)
		return nil, err
	}
	s.logger.Info("listed tags", "registry_id", registryID, "repository", repoName, "count", len(tags), "duration", duration)
	return tags, nil
}

// DeleteRegistry deletes the registry.
func (s *Service) DeleteRegistry(ctx context.Context, token string, registryID string) error {
	s.logger.Info("deleting registry", "registry_id", registryID)
	client := v1.NewCRaaSClientV1(token, s.endpoint)

	start := time.Now()
	_, err := registry.Delete(ctx, client, registryID)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to delete registry", "registry_id", registryID, "error", err)
		return err
	}
	s.logger.Info("registry deleted", "registry_id", registryID, "duration", duration)
	return nil
}

// DeleteRepository deletes the repository.
func (s *Service) DeleteRepository(ctx context.Context, token string, registryID, repoName string) error {
	s.logger.Info("deleting repository", "registry_id", registryID, "repository", repoName)
	client := v1.NewCRaaSClientV1(token, s.endpoint)

	start := time.Now()
	_, err := repository.DeleteRepository(ctx, client, registryID, repoName)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to delete repository", "registry_id", registryID, "repository", repoName, "error", err)
		return err
	}
	s.logger.Info("repository deleted", "registry_id", registryID, "repository", repoName, "duration", duration)
	return nil
}

// DeleteImage deletes the image by digest.
func (s *Service) DeleteImage(ctx context.Context, token string, registryID, repoName, digest string) error {
	s.logger.Info("deleting image", "registry_id", registryID, "repository", repoName, "digest", digest)
	client := v1.NewCRaaSClientV1(token, s.endpoint)

	start := time.Now()
	_, err := repository.DeleteImageManifest(ctx, client, registryID, repoName, digest)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to delete image", "registry_id", registryID, "repository", repoName, "digest", digest, "error", err)
		return err
	}
	s.logger.Info("image deleted", "registry_id", registryID, "repository", repoName, "digest", digest, "duration", duration)
	return nil
}

// CleanupResult represents the result of a cleanup operation.
type CleanupResult struct {
	Deleted []interface{} `json:"deleted"`
	Failed  []interface{} `json:"failed"`
}

// CleanupRequest represents the request body for cleanup operation.
type CleanupRequest struct {
	Digests   []string `json:"digests"`
	DisableGC bool     `json:"disable_gc"`
	Tags      []string `json:"tags,omitempty"`
}

// CleanupRepository cleans up the repository.
func (s *Service) CleanupRepository(ctx context.Context, token, registryID, repoName string, digests []string, disableGC bool) (*CleanupResult, error) {
	s.logger.Info("cleaning up repository", "registry_id", registryID, "repository", repoName, "digest_count", len(digests), "disable_gc", disableGC)

	url := fmt.Sprintf("%s/registries/%s/repositories/%s/cleanup", s.endpoint, registryID, repoName)

	reqBody := CleanupRequest{
		Digests:   digests,
		DisableGC: disableGC,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal cleanup request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create cleanup request: %w", err)
	}

	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("Content-Type", "application/json")

	// Use a client with appropriate timeout
	httpClient := &http.Client{
		Timeout: 60 * time.Second, // Increased timeout for potentially long-running cleanup
	}

	start := time.Now()
	resp, err := httpClient.Do(req)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to execute cleanup request", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("cleanup request failed", "status", resp.StatusCode, "body", string(body))
		return nil, fmt.Errorf("cleanup request failed with status: %d", resp.StatusCode)
	}

	var result CleanupResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		s.logger.Error("failed to decode cleanup response", "error", err)
		return nil, fmt.Errorf("failed to decode cleanup response: %w", err)
	}

	s.logger.Info("repository cleanup completed", "registry_id", registryID, "repository", repoName, "duration", duration, "deleted_count", len(result.Deleted), "failed_count", len(result.Failed))
	return &result, nil
}
