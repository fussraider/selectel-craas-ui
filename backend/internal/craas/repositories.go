package craas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	clientv1 "github.com/selectel/craas-go/pkg/v1/client"
	"github.com/selectel/craas-go/pkg/v1/repository"
)

// ListRepositories returns a list of repositories in the registry.
func (s *Service) ListRepositories(ctx context.Context, token string, registryID string) ([]*repository.Repository, error) {
	s.logger.Debug("listing repositories", "registry_id", registryID)
	client, err := clientv1.NewCRaaSClientV1(token, s.endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

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

// DeleteRepository deletes the repository.
func (s *Service) DeleteRepository(ctx context.Context, token string, registryID, repoName string) error {
	s.logger.Info("deleting repository", "registry_id", registryID, "repository", repoName)
	client, err := clientv1.NewCRaaSClientV1(token, s.endpoint)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	encodedRepoName := url.PathEscape(repoName)

	start := time.Now()
	_, err = repository.DeleteRepository(ctx, client, registryID, encodedRepoName)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to delete repository", "registry_id", registryID, "repository", repoName, "error", err)
		return err
	}
	s.logger.Info("repository deleted", "registry_id", registryID, "repository", repoName, "duration", duration)
	return nil
}

// CleanupRepository cleans up the repository.
func (s *Service) CleanupRepository(ctx context.Context, token, registryID, repoName string, digests []string, disableGC bool) (*CleanupResult, error) {
	s.logger.Info("cleaning up repository", "registry_id", registryID, "repository", repoName, "digest_count", len(digests), "disable_gc", disableGC)

	encodedRepoName := url.PathEscape(repoName)
	encodedRegistryID := url.PathEscape(registryID)
	cleanupUrl := fmt.Sprintf("%s/registries/%s/repositories/%s/cleanup", s.endpoint, encodedRegistryID, encodedRepoName)

	reqBody := CleanupRequest{
		Digests:   digests,
		DisableGC: disableGC,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal cleanup request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", cleanupUrl, bytes.NewBuffer(bodyBytes))
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
