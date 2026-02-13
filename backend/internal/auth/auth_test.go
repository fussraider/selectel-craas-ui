package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/generic/selectel-craas-web/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetAccountToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v3/auth/tokens" && r.Method == "POST" {
			w.Header().Set("X-Subject-Token", "fake-token")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"token": {"expires_at": "2030-01-01T00:00:00Z"}}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	cfg := &config.Config{
		SelectelUsername:  "testuser",
		SelectelAccountID: "12345",
		SelectelPassword:  "password",
	}
	client := New(cfg)
	client.AuthURL = ts.URL + "/v3/auth/tokens"

	token, err := client.GetAccountToken()
	assert.NoError(t, err)
	assert.Equal(t, "fake-token", token)
}

func TestListProjects(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v3/auth/projects" && r.Method == "GET" {
			token := r.Header.Get("X-Auth-Token")
			if token != "fake-token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"projects": [{"id": "p1", "name": "Project 1"}]}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	cfg := &config.Config{}
	client := New(cfg)
	client.ProjURL = ts.URL + "/v3/auth/projects"

	projects, err := client.ListProjects("fake-token")
	assert.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, "p1", projects[0].ID)
}

func TestGetProjectToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v3/auth/tokens" && r.Method == "POST" {
			var body map[string]interface{}
			json.NewDecoder(r.Body).Decode(&body)
			scope := body["auth"].(map[string]interface{})["scope"].(map[string]interface{})
			proj := scope["project"].(map[string]interface{})

			if proj["id"] == "p1" {
				w.Header().Set("X-Subject-Token", "project-token")
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{}`))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	cfg := &config.Config{
		SelectelUsername:  "testuser",
		SelectelAccountID: "12345",
		SelectelPassword:  "password",
	}
	client := New(cfg)
	client.AuthURL = ts.URL + "/v3/auth/tokens"

	token, err := client.GetProjectToken("p1")
	assert.NoError(t, err)
	assert.Equal(t, "project-token", token)
}
