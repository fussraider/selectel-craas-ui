package craas

import (
	"context"
	"log/slog"
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
	s.logger.Info("listed images", "registry_id", registryID, "repository", repoName, "count", len(images), "duration", duration)
	return images, nil
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
