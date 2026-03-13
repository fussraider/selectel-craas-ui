package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/generic/selectel-craas-web/internal/auth"
	"github.com/generic/selectel-craas-web/internal/config"
	"github.com/generic/selectel-craas-web/internal/craas"
	"github.com/golang-jwt/jwt/v5"
)

var testLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func TestEnableCORS(t *testing.T) {
	tests := []struct {
		name           string
		allowedOrigin  string
		requestMethod  string
		expectedOrigin string
		expectedStatus int
	}{
		{
			name:           "Default Origin (*)",
			allowedOrigin:  "*",
			requestMethod:  "GET",
			expectedOrigin: "*",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Specific Origin",
			allowedOrigin:  "https://example.com",
			requestMethod:  "GET",
			expectedOrigin: "https://example.com",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "OPTIONS Request",
			allowedOrigin:  "*",
			requestMethod:  "OPTIONS",
			expectedOrigin: "*",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			cfg := &config.Config{
				CORSAllowedOrigin: tt.allowedOrigin,
			}
			server := &Server{
				Config: cfg,
			}

			// Create a dummy handler that the middleware will wrap
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create the middleware
			handler := server.EnableCORS(nextHandler)

			// Create a request
			req := httptest.NewRequest(tt.requestMethod, "/test", nil)
			rr := httptest.NewRecorder()

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			// Check the Access-Control-Allow-Origin header
			if origin := rr.Header().Get("Access-Control-Allow-Origin"); origin != tt.expectedOrigin {
				t.Errorf("handler returned wrong Access-Control-Allow-Origin header: got %v want %v",
					origin, tt.expectedOrigin)
			}

			// Check the Access-Control-Allow-Credentials header
			if credentials := rr.Header().Get("Access-Control-Allow-Credentials"); credentials != "true" {
				t.Errorf("handler returned wrong Access-Control-Allow-Credentials header: got %v want %v",
					credentials, "true")
			}

			// Check other CORS headers
			if methods := rr.Header().Get("Access-Control-Allow-Methods"); methods == "" {
				t.Errorf("handler returned empty Access-Control-Allow-Methods header")
			}
			if headers := rr.Header().Get("Access-Control-Allow-Headers"); headers == "" {
				t.Errorf("handler returned empty Access-Control-Allow-Headers header")
			}
		})
	}
}

func TestExecuteWithRetry(t *testing.T) {
	// A helper to quickly set up a mock auth server.
	setupMockAuth := func(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *auth.Client) {
		ts := httptest.NewServer(handler)
		t.Cleanup(ts.Close)

		cfg := &config.Config{
			SelectelUsername:  "testuser",
			SelectelAccountID: "12345",
			SelectelPassword:  "password",
		}
		client := auth.New(cfg, testLogger)
		client.AuthURL = ts.URL + "/v3/auth/tokens"
		return ts, client
	}

	t.Run("Success on first attempt", func(t *testing.T) {
		var tokenFetchCount int
		ts, authClient := setupMockAuth(t, func(w http.ResponseWriter, r *http.Request) {
			tokenFetchCount++
			w.Header().Set("X-Subject-Token", "token1")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{}`))
		})
		_ = ts

		server := &Server{
			Auth:   authClient,
			Logger: testLogger,
		}

		var opCalls int
		err := server.ExecuteWithRetry(context.Background(), "pid1", func(token string) error {
			opCalls++
			if token != "token1" {
				return fmt.Errorf("expected token1, got %s", token)
			}
			return nil
		})

		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if tokenFetchCount != 1 {
			t.Errorf("expected 1 token fetch, got %d", tokenFetchCount)
		}
		if opCalls != 1 {
			t.Errorf("expected 1 op call, got %d", opCalls)
		}
	})

	t.Run("Initial token fetch fails, retry succeeds", func(t *testing.T) {
		var tokenFetchCount int
		ts, authClient := setupMockAuth(t, func(w http.ResponseWriter, r *http.Request) {
			tokenFetchCount++
			if tokenFetchCount == 1 {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("X-Subject-Token", "token2")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{}`))
		})
		_ = ts

		server := &Server{
			Auth:   authClient,
			Logger: testLogger,
		}

		var opCalls int
		err := server.ExecuteWithRetry(context.Background(), "pid2", func(token string) error {
			opCalls++
			if token != "token2" {
				return fmt.Errorf("expected token2, got %s", token)
			}
			return nil
		})

		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if tokenFetchCount != 2 {
			t.Errorf("expected 2 token fetches, got %d", tokenFetchCount)
		}
		if opCalls != 1 {
			t.Errorf("expected 1 op call, got %d", opCalls)
		}
	})

	t.Run("Operation returns standard error, fails immediately without retry", func(t *testing.T) {
		ts, authClient := setupMockAuth(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Subject-Token", "token1")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{}`))
		})
		_ = ts

		server := &Server{
			Auth:   authClient,
			Logger: testLogger,
		}

		var opCalls int
		expectedErr := errors.New("standard error")
		err := server.ExecuteWithRetry(context.Background(), "pid3", func(token string) error {
			opCalls++
			return expectedErr
		})

		if err != expectedErr {
			t.Fatalf("expected error %v, got %v", expectedErr, err)
		}
		if opCalls != 1 {
			t.Errorf("expected 1 op call, got %d", opCalls)
		}
	})

	t.Run("Operation returns Auth error, succeeds on retry", func(t *testing.T) {
		var tokenFetchCount int
		ts, authClient := setupMockAuth(t, func(w http.ResponseWriter, r *http.Request) {
			tokenFetchCount++
			w.Header().Set("X-Subject-Token", fmt.Sprintf("token%d", tokenFetchCount))
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{}`))
		})
		_ = ts

		server := &Server{
			Auth:   authClient,
			Logger: testLogger,
		}

		var opCalls int
		err := server.ExecuteWithRetry(context.Background(), "pid4", func(token string) error {
			opCalls++
			if opCalls == 1 {
				if token != "token1" {
					return fmt.Errorf("expected token1, got %s", token)
				}
				return craas.ErrUnauthorized
			}
			if token != "token2" {
				return fmt.Errorf("expected token2, got %s", token)
			}
			return nil
		})

		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if tokenFetchCount != 2 {
			t.Errorf("expected 2 token fetches, got %d", tokenFetchCount)
		}
		if opCalls != 2 {
			t.Errorf("expected 2 op calls, got %d", opCalls)
		}
	})

	t.Run("Operation returns Auth error, token re-fetch fails", func(t *testing.T) {
		var tokenFetchCount int
		ts, authClient := setupMockAuth(t, func(w http.ResponseWriter, r *http.Request) {
			tokenFetchCount++
			if tokenFetchCount == 1 {
				w.Header().Set("X-Subject-Token", "token1")
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{}`))
				return
			}
			// Second fetch fails
			w.WriteHeader(http.StatusInternalServerError)
		})
		_ = ts

		server := &Server{
			Auth:   authClient,
			Logger: testLogger,
		}

		var opCalls int
		err := server.ExecuteWithRetry(context.Background(), "pid5", func(token string) error {
			opCalls++
			return craas.ErrUnauthorized
		})

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if tokenFetchCount != 2 {
			t.Errorf("expected 2 token fetches, got %d", tokenFetchCount)
		}
		if opCalls != 1 {
			t.Errorf("expected 1 op call, got %d", opCalls)
		}
	})
}

func TestSecurityHeaders(t *testing.T) {
	server := &Server{}

	// Create a dummy handler that the middleware will wrap
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create the middleware
	handler := server.SecurityHeaders(nextHandler)

	// Create a request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check headers
	headers := map[string]string{
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "DENY",
		"Content-Security-Policy": "default-src 'self'",
		"Referrer-Policy":         "strict-origin-when-cross-origin",
	}

	for k, v := range headers {
		if val := rr.Header().Get(k); val != v {
			t.Errorf("header %s: got %s want %s", k, val, v)
		}
	}
}

func TestAuthMiddleware(t *testing.T) {
	secret := "test-secret"
	cfg := &config.Config{
		AuthEnabled: true,
		JWTSecret:   secret,
	}
	server := &Server{
		Config: cfg,
	}

	// Create a valid token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "testuser",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(string)
		if !ok || user != "testuser" {
			w.WriteHeader(http.StatusForbidden) // User not in context or wrong
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := server.AuthMiddleware(nextHandler)

	tests := []struct {
		name           string
		cookieName     string
		cookieValue    string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "No Auth",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Valid Cookie",
			cookieName:     "auth_token",
			cookieValue:    tokenString,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid Header",
			authHeader:     "Bearer " + tokenString,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Cookie",
			cookieName:     "auth_token",
			cookieValue:    "invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Header",
			authHeader:     "Bearer invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Wrong Cookie Name",
			cookieName:     "wrong_token",
			cookieValue:    tokenString,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.cookieName != "" {
				req.AddCookie(&http.Cookie{Name: tt.cookieName, Value: tt.cookieValue})
			}
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
