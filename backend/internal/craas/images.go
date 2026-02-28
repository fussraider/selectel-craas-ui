package craas

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	clientv1 "github.com/selectel/craas-go/pkg/v1/client"
	"github.com/selectel/craas-go/pkg/v1/repository"
)

// ListImages returns a list of images in the repository.
func (s *Service) ListImages(ctx context.Context, token string, registryID, repoName string) ([]*repository.Image, error) {
	s.logger.Debug("listing images", "registry_id", registryID, "repository", repoName)
	client, err := clientv1.NewCRaaSClientV1(token, s.endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	// URL-encode repository name to handle special characters like slashes
	encodedRepoName := url.PathEscape(repoName)

	start := time.Now()
	images, _, err := repository.ListImages(ctx, client, registryID, encodedRepoName)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to list images", "registry_id", registryID, "repository", repoName, "error", err)
		return nil, err
	}

	if s.enableMissingTagsCheck {
		images = s.resolveMissingTags(ctx, client, token, registryID, repoName, encodedRepoName, images)
	}

	s.logger.Info("listed images", "registry_id", registryID, "repository", repoName, "count", len(images), "duration", duration)
	return images, nil
}

func (s *Service) resolveMissingTags(ctx context.Context, client *clientv1.ServiceClient, token, registryID, repoName, encodedRepoName string, images []*repository.Image) []*repository.Image {
	// Fetch all tags to check for missing ones
	allTags, _, err := repository.ListTags(ctx, client, registryID, encodedRepoName)
	if err != nil {
		s.logger.Warn("failed to list tags for verification", "registry_id", registryID, "repository", repoName, "error", err)
		return images
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

	if len(missingTags) == 0 {
		return images
	}

	s.logger.Info("found missing tags, resolving", "count", len(missingTags), "tags", missingTags)

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(5) // Limit concurrency

	var mu sync.Mutex

	for _, tag := range missingTags {
		tag := tag
		g.Go(func() error {
			// Fetch all digests associated with the tag
			// Use encodedRepoName here too
			digests, err := s.fetchImageDigests(ctx, httpClient, token, registryID, encodedRepoName, tag)
			if err != nil {
				s.logger.Warn("failed to fetch digests for tag", "tag", tag, "error", err)
				return nil // Don't fail the whole group, just skip
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
			return nil
		})
	}
	// Wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		s.logger.Error("error resolving missing tags", "error", err)
	}

	return images
}

// fetchImageDigests fetches the digest(s) associated with a tag.
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
		}
	}

	if len(digests) == 0 {
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
	client, err := clientv1.NewCRaaSClientV1(token, s.endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	encodedRepoName := url.PathEscape(repoName)

	start := time.Now()
	tags, _, err := repository.ListTags(ctx, client, registryID, encodedRepoName)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to list tags", "registry_id", registryID, "repository", repoName, "error", err)
		return nil, err
	}
	s.logger.Info("listed tags", "registry_id", registryID, "repository", repoName, "count", len(tags), "duration", duration)
	return tags, nil
}

// DeleteImage deletes the image by digest.
func (s *Service) DeleteImage(ctx context.Context, token string, registryID, repoName, digest string) error {
	s.logger.Info("deleting image", "registry_id", registryID, "repository", repoName, "digest", digest)
	client, err := clientv1.NewCRaaSClientV1(token, s.endpoint)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	encodedRepoName := url.PathEscape(repoName)

	start := time.Now()
	_, err = repository.DeleteImageManifest(ctx, client, registryID, encodedRepoName, digest)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to delete image", "registry_id", registryID, "repository", repoName, "digest", digest, "error", err)
		return err
	}
	s.logger.Info("image deleted", "registry_id", registryID, "repository", repoName, "digest", digest, "duration", duration)
	return nil
}
