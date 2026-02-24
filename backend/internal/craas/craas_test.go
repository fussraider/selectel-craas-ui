package craas

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/selectel/craas-go/pkg/v1/repository"
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
	digest := "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/registries/reg1/repositories/repo1/images" && r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"digest": "` + digest + `"}]`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	svc := &Service{endpoint: ts.URL + "/v1", logger: testLogger}
	images, err := svc.ListImages(context.Background(), "fake-token", "reg1", "repo1")

	assert.NoError(t, err)
	assert.Len(t, images, 1)
	assert.Equal(t, digest, images[0].Digest)
}

func TestListImages_EncodedRepoName(t *testing.T) {
	digest := "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strictly check for encoded slash
		if strings.Contains(r.RequestURI, "group%2Frepo") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"digest": "` + digest + `"}]`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	svc := &Service{endpoint: ts.URL + "/v1", logger: testLogger}
	// We pass "group/repo" as the name. It should be encoded by the client.
	images, err := svc.ListImages(context.Background(), "fake-token", "reg1", "group/repo")

	assert.NoError(t, err)
	assert.Len(t, images, 1)
	assert.Equal(t, digest, images[0].Digest)
}

func TestListImages_MissingTags(t *testing.T) {
	// 64-char hashes
	hAbc := "sha256:a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1"
	hDef := "sha256:b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2"
	hGhi := "sha256:c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3"
	hSub1 := "sha256:d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4"
	hSub2 := "sha256:e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5"
	hNested := "sha256:f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6"
	hAccept := "sha256:0000000000000000000000000000000000000000000000000000000000000000"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/registries/reg1/repositories/repo1/images":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// Returns images
			w.Write([]byte(`[
				{"digest": "` + hAbc + `", "tags": []},
				{"digest": "` + hSub1 + `", "tags": []},
				{"digest": "` + hSub2 + `", "tags": []},
                {"digest": "` + hNested + `", "tags": []},
                {"digest": "` + hAccept + `", "tags": []}
			]`))
		case "/v1/registries/reg1/repositories/repo1/tags":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// Returns tags including "v6-accept"
			w.Write([]byte(`["v1", "v2", "v3", "v4-list", "v5-nested", "v6-accept"]`))
		case "/v1/registries/reg1/repositories/repo1/v1":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"digest": "` + hAbc + `", "tags": ["v1"], "size": 100}`))
		case "/v1/registries/reg1/repositories/repo1/v2":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"digest": "` + hDef + `", "tags": ["v2"], "size": 200}`))
		case "/v1/registries/reg1/repositories/repo1/v3":
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Docker-Content-Digest", hGhi)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"tags": ["v3"], "size": 300}`))
		case "/v1/registries/reg1/repositories/repo1/v4-list":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"digest": "` + hSub1 + `", "size": 100}, {"digest": "` + hSub2 + `", "size": 200}]`))
		case "/v1/registries/reg1/repositories/repo1/v5-nested":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"foo": [{"bar": {"baz": "` + hNested + `"}}]}`))
		case "/v1/registries/reg1/repositories/repo1/v6-accept":
			// Check for Accept header
			accept := r.Header.Get("Accept")
			if strings.Contains(accept, "application/vnd.docker.distribution.manifest.list.v2+json") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				// Return valid JSON
				w.Write([]byte(`{"digest": "` + hAccept + `", "tags": ["v6-accept"]}`))
			} else {
				// Simulate failure (empty array) without header
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[]`))
			}

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	svc := &Service{endpoint: ts.URL + "/v1", logger: testLogger}
	images, err := svc.ListImages(context.Background(), "fake-token", "reg1", "repo1")

	assert.NoError(t, err)

	// Check if we correctly handle images
	var imgAccept *repository.Image

	for _, img := range images {
		if img.Digest == hAccept {
			imgAccept = img
		}
	}

	assert.NotNil(t, imgAccept)
	assert.Contains(t, imgAccept.Tags, "v6-accept")
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

func TestCleanupRepository(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Expect encoded path
		if strings.Contains(r.RequestURI, "group%2Frepo/cleanup") {
			var req CleanupRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if len(req.Digests) == 1 && req.Digests[0] == "sha256:abc" && req.DisableGC == false {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"deleted": [{"digest": "sha256:abc", "tags": []}],
					"failed": []
				}`))
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	svc := &Service{endpoint: ts.URL + "/v1", logger: testLogger}
	result, err := svc.CleanupRepository(context.Background(), "fake-token", "reg1", "group/repo", []string{"sha256:abc"}, false)

	assert.NoError(t, err)
	assert.Len(t, result.Deleted, 1)
	assert.Len(t, result.Failed, 0)
}
