package craas

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func TestListRegistries(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/registries" && r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id": "reg1", "name": "registry-1"}]`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	svc := &Service{endpoint: ts.URL + "/v1", logger: testLogger}
	registries, err := svc.ListRegistries(context.Background(), "fake-token")

	assert.NoError(t, err)
	assert.Len(t, registries, 1)
	assert.Equal(t, "reg1", registries[0].ID)
}

func TestListRepositories(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/registries/reg1/repositories" && r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"name": "repo1"}]`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	svc := &Service{endpoint: ts.URL + "/v1", logger: testLogger}
	repos, err := svc.ListRepositories(context.Background(), "fake-token", "reg1")

	assert.NoError(t, err)
	assert.Len(t, repos, 1)
	assert.Equal(t, "repo1", repos[0].Name)
}

func TestListImages(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/registries/reg1/repositories/repo1/images" && r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"digest": "sha256:abc"}]`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	svc := &Service{endpoint: ts.URL + "/v1", logger: testLogger}
	images, err := svc.ListImages(context.Background(), "fake-token", "reg1", "repo1")

	assert.NoError(t, err)
	assert.Len(t, images, 1)
	assert.Equal(t, "sha256:abc", images[0].Digest)
}

func TestDeleteRegistry(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/registries/reg1" && r.Method == "DELETE" {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	svc := &Service{endpoint: ts.URL + "/v1", logger: testLogger}
	err := svc.DeleteRegistry(context.Background(), "fake-token", "reg1")
	assert.NoError(t, err)
}
