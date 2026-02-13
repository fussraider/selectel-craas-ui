package craas

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
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

				img, err := s.fetchImageManifest(ctx, httpClient, token, registryID, repoName, tag)
				if err != nil {
					s.logger.Warn("failed to fetch manifest for tag", "tag", tag, "error", err)
					return
				}

				mu.Lock()
				defer mu.Unlock()

				// Check if we already have this image (by digest)
				if existingImg, ok := imageMap[img.Digest]; ok {
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
					// New image found (e.g. manifest list)
					// Ensure the image has the tag
					hasTag := false
					for _, t := range img.Tags {
						if t == tag {
							hasTag = true
							break
						}
					}
					if !hasTag {
						img.Tags = append(img.Tags, tag)
					}

					images = append(images, img)
					imageMap[img.Digest] = img
				}
			}(tag)
		}
		wg.Wait()
	}

	s.logger.Info("listed images", "registry_id", registryID, "repository", repoName, "count", len(images), "duration", duration)
	return images, nil
}

// fetchImageManifest fetches the image manifest details for a given tag/digest.
func (s *Service) fetchImageManifest(ctx context.Context, client *http.Client, token, registryID, repoName, reference string) (*repository.Image, error) {
	url := fmt.Sprintf("%s/registries/%s/repositories/%s/%s", s.endpoint, registryID, repoName, reference)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var img repository.Image
	if err := json.NewDecoder(resp.Body).Decode(&img); err != nil {
		return nil, err
	}

	// Fallback to header if body digest is empty
	if img.Digest == "" {
		img.Digest = resp.Header.Get("Docker-Content-Digest")
	}

	return &img, nil
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
