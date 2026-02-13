package craas

import (
	"context"

	v1 "github.com/selectel/craas-go/pkg"
	"github.com/selectel/craas-go/pkg/v1/registry"
	"github.com/selectel/craas-go/pkg/v1/repository"
)

const Endpoint = "https://cr.selcloud.ru/api/v1"

type Service struct {
	endpoint string
}

func New() *Service {
	return &Service{endpoint: Endpoint}
}

// ListRegistries returns a list of registries for the project (scoped by token).
func (s *Service) ListRegistries(ctx context.Context, token string) ([]*registry.Registry, error) {
	client := v1.NewCRaaSClientV1(token, s.endpoint)
	registries, _, err := registry.List(ctx, client)
	if err != nil {
		return nil, err
	}
	return registries, nil
}

// ListRepositories returns a list of repositories in the registry.
func (s *Service) ListRepositories(ctx context.Context, token string, registryID string) ([]*repository.Repository, error) {
	client := v1.NewCRaaSClientV1(token, s.endpoint)
	repos, _, err := repository.ListRepositories(ctx, client, registryID)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// ListImages returns a list of images in the repository.
func (s *Service) ListImages(ctx context.Context, token string, registryID, repoName string) ([]*repository.Image, error) {
	client := v1.NewCRaaSClientV1(token, s.endpoint)
	images, _, err := repository.ListImages(ctx, client, registryID, repoName)
	if err != nil {
		return nil, err
	}
	return images, nil
}

// ListTags returns a list of tags in the repository.
func (s *Service) ListTags(ctx context.Context, token string, registryID, repoName string) ([]string, error) {
	client := v1.NewCRaaSClientV1(token, s.endpoint)
	tags, _, err := repository.ListTags(ctx, client, registryID, repoName)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// DeleteRegistry deletes the registry.
func (s *Service) DeleteRegistry(ctx context.Context, token string, registryID string) error {
	client := v1.NewCRaaSClientV1(token, s.endpoint)
	_, err := registry.Delete(ctx, client, registryID)
	return err
}

// DeleteRepository deletes the repository.
func (s *Service) DeleteRepository(ctx context.Context, token string, registryID, repoName string) error {
	client := v1.NewCRaaSClientV1(token, s.endpoint)
	_, err := repository.DeleteRepository(ctx, client, registryID, repoName)
	return err
}

// DeleteImage deletes the image by digest.
func (s *Service) DeleteImage(ctx context.Context, token string, registryID, repoName, digest string) error {
	client := v1.NewCRaaSClientV1(token, s.endpoint)
	_, err := repository.DeleteImageManifest(ctx, client, registryID, repoName, digest)
	return err
}
