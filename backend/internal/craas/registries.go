package craas

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	clientv1 "github.com/selectel/craas-go/pkg/v1/client"
	"github.com/selectel/craas-go/pkg/v1/registry"
)

// ListRegistries returns a list of registries for the project (scoped by token).
func (s *Service) ListRegistries(ctx context.Context, token string) ([]*registry.Registry, error) {
	s.logger.Debug("listing registries")
	client, err := clientv1.NewCRaaSClientV1(token, s.endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

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

// DeleteRegistry deletes the registry.
func (s *Service) DeleteRegistry(ctx context.Context, token string, registryID string) error {
	s.logger.Info("deleting registry", "registry_id", registryID)
	client, err := clientv1.NewCRaaSClientV1(token, s.endpoint)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	start := time.Now()
	_, err = registry.Delete(ctx, client, registryID)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to delete registry", "registry_id", registryID, "error", err)
		return err
	}
	s.logger.Info("registry deleted", "registry_id", registryID, "duration", duration)
	return nil
}

// GetGCInfo returns the garbage collection size information.
func (s *Service) GetGCInfo(ctx context.Context, token, registryID string) (*GCInfo, error) {
	s.logger.Debug("getting gc info", "registry_id", registryID)

	url := fmt.Sprintf("%s/registries/%s/garbage-collection/size", s.endpoint, registryID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Auth-Token", token)

	httpClient := &http.Client{
		Timeout: 5 * time.Minute,
	}

	start := time.Now()
	resp, err := httpClient.Do(req)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to get gc info", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnauthorized
		}
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("get gc info failed", "status", resp.StatusCode, "body", string(body))
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	var info GCInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		s.logger.Error("failed to decode gc info", "error", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	s.logger.Info("got gc info", "registry_id", registryID, "duration", duration)
	return &info, nil
}

// StartGC initiates the garbage collection process.
func (s *Service) StartGC(ctx context.Context, token, registryID string) error {
	s.logger.Info("starting gc", "registry_id", registryID)

	url := fmt.Sprintf("%s/registries/%s/garbage-collection", s.endpoint, registryID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Auth-Token", token)

	httpClient := &http.Client{
		Timeout: 5 * time.Minute,
	}

	start := time.Now()
	resp, err := httpClient.Do(req)
	duration := time.Since(start)

	if err != nil {
		s.logger.Error("failed to start gc", "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusUnauthorized {
			return ErrUnauthorized
		}
		if resp.StatusCode == http.StatusConflict {
			return fmt.Errorf("garbage collection already in progress")
		}
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("start gc failed", "status", resp.StatusCode, "body", string(body))
		return fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	s.logger.Info("gc started", "registry_id", registryID, "duration", duration)
	return nil
}
