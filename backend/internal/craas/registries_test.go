package craas

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/generic/selectel-craas-web/internal/config"
)

func TestService_GetGCInfo(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/registries/reg-1/garbage-collection/size" {
			t.Errorf("expected path /registries/reg-1/garbage-collection/size, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.Header.Get("X-Auth-Token") != "test-token" {
			t.Errorf("expected test-token, got %s", r.Header.Get("X-Auth-Token"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GCInfo{
			SizeNonReferenced: 100,
			SizeSummary:       200,
			SizeUntagged:      50,
		})
	}))
	defer server.Close()

	// Service
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc := New(&config.Config{}, logger)
	svc.endpoint = server.URL // Override endpoint

	// Test
	info, err := svc.GetGCInfo(context.Background(), "test-token", "reg-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.SizeSummary != 200 {
		t.Errorf("expected size 200, got %d", info.SizeSummary)
	}
}

func TestService_StartGC(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/registries/reg-1/garbage-collection" {
			t.Errorf("expected path /registries/reg-1/garbage-collection, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	// Service
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc := New(&config.Config{}, logger)
	svc.endpoint = server.URL

	// Test
	err := svc.StartGC(context.Background(), "test-token", "reg-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_StartGC_Conflict(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
	}))
	defer server.Close()

	// Service
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc := New(&config.Config{}, logger)
	svc.endpoint = server.URL

	// Test
	err := svc.StartGC(context.Background(), "test-token", "reg-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "garbage collection already in progress" {
		t.Errorf("expected 'garbage collection already in progress', got '%s'", err.Error())
	}
}

func TestService_StartGC_Unauthorized(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	// Service
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc := New(&config.Config{}, logger)
	svc.endpoint = server.URL

	// Test
	err := svc.StartGC(context.Background(), "test-token", "reg-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestService_StartGC_GenericError(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Service
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc := New(&config.Config{}, logger)
	svc.endpoint = server.URL

	// Test
	err := svc.StartGC(context.Background(), "test-token", "reg-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expected := "request failed with status: 500"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}
