package craas

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestListImages_MissingTags(t *testing.T) {
	// 64-char hashes
	hAbc := "sha256:a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1"
	hDef := "sha256:b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2"
	hGhi := "sha256:c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3"
	hSub1 := "sha256:d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4"
	hSub2 := "sha256:e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5e5"
	hNested := "sha256:f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/registries/reg1/repositories/repo1/images":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// Returns one image "abc" without tags
			w.Write([]byte(`[
				{"digest": "` + hAbc + `", "tags": []},
				{"digest": "` + hSub1 + `", "tags": []},
				{"digest": "` + hSub2 + `", "tags": []},
                {"digest": "` + hNested + `", "tags": []}
			]`))
		case "/v1/registries/reg1/repositories/repo1/tags":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// Returns "v1", "v2", "v3", "v4-list", "v5-nested"
			w.Write([]byte(`["v1", "v2", "v3", "v4-list", "v5-nested"]`))
		case "/v1/registries/reg1/repositories/repo1/v1":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// v1 points to abc
			w.Write([]byte(`{"digest": "` + hAbc + `", "tags": ["v1"], "size": 100}`))
		case "/v1/registries/reg1/repositories/repo1/v2":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// v2 points to def (new image)
			w.Write([]byte(`{"digest": "` + hDef + `", "tags": ["v2"], "size": 200}`))
		case "/v1/registries/reg1/repositories/repo1/v3":
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Docker-Content-Digest", hGhi)
			w.WriteHeader(http.StatusOK)
			// v3 points to ghi (new image), missing body digest
			w.Write([]byte(`{"tags": ["v3"], "size": 300}`))
		case "/v1/registries/reg1/repositories/repo1/v4-list":
			w.Header().Set("Content-Type", "application/json")
			// NO Docker-Content-Digest header
			w.WriteHeader(http.StatusOK)
			// v4-list is a manifest list returned as an array
			// containing digests that match "sub1" and "sub2" in ListImages
			w.Write([]byte(`[{"digest": "` + hSub1 + `", "size": 100}, {"digest": "` + hSub2 + `", "size": 200}]`))
		case "/v1/registries/reg1/repositories/repo1/v5-nested":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// v5-nested is a deeply nested JSON with "sha256:nested" hidden in it
			w.Write([]byte(`{"foo": [{"bar": {"baz": "` + hNested + `"}}]}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	svc := &Service{endpoint: ts.URL + "/v1", logger: testLogger}
	images, err := svc.ListImages(context.Background(), "fake-token", "reg1", "repo1")

	assert.NoError(t, err)
	// We expect images: abc, sub1, sub2, nested, def, ghi

	// Check if we correctly handle images
	var imgAbc *repository.Image
	var imgSub1 *repository.Image
	var imgSub2 *repository.Image
	var imgNested *repository.Image
	var imgDef *repository.Image
	var imgGhi *repository.Image

	for _, img := range images {
		if img.Digest == hAbc {
			imgAbc = img
		} else if img.Digest == hSub1 {
			imgSub1 = img
		} else if img.Digest == hSub2 {
			imgSub2 = img
		} else if img.Digest == hNested {
			imgNested = img
		} else if img.Digest == hDef {
			imgDef = img
		} else if img.Digest == hGhi {
			imgGhi = img
		}
	}

	assert.NotNil(t, imgAbc)
	assert.Contains(t, imgAbc.Tags, "v1")

	assert.NotNil(t, imgSub1)
	assert.Contains(t, imgSub1.Tags, "v4-list")

	assert.NotNil(t, imgSub2)
	assert.Contains(t, imgSub2.Tags, "v4-list")

    assert.NotNil(t, imgNested)
    assert.Contains(t, imgNested.Tags, "v5-nested")

	assert.NotNil(t, imgDef)
	assert.Contains(t, imgDef.Tags, "v2")

	assert.NotNil(t, imgGhi)
	assert.Contains(t, imgGhi.Tags, "v3")
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
