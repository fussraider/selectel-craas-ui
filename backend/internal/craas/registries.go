package craas

import (
	"context"
	"time"

	v1 "github.com/selectel/craas-go/pkg"
	"github.com/selectel/craas-go/pkg/v1/registry"
)

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
