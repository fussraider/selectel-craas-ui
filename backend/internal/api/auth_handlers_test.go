package api

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/generic/selectel-craas-web/internal/config"
)

func TestLogin(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	tests := []struct {
		name           string
		authEnabled    bool
		authLogin      string
		authPassword   string
		jwtSecret      string
		requestBody    interface{}
		expectedStatus int
		expectToken    bool
	}{
		{
			name:           "Auth Disabled",
			authEnabled:    false,
			requestBody:    LoginRequest{Login: "user", Password: "password"},
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
		},
		{
			name:           "Invalid JSON",
			authEnabled:    true,
			requestBody:    "invalid-json",
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
		},
		{
			name:           "Invalid Password",
			authEnabled:    true,
			authLogin:      "admin",
			authPassword:   "password",
			requestBody:    LoginRequest{Login: "admin", Password: "wrong-password"},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
		{
			name:           "Invalid Login",
			authEnabled:    true,
			authLogin:      "admin",
			authPassword:   "password",
			requestBody:    LoginRequest{Login: "wrong-user", Password: "password"},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
		{
			name:           "Successful Login",
			authEnabled:    true,
			authLogin:      "admin",
			authPassword:   "password",
			jwtSecret:      "secret",
			requestBody:    LoginRequest{Login: "admin", Password: "password"},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				AuthEnabled:  tt.authEnabled,
				AuthLogin:    tt.authLogin,
				AuthPassword: tt.authPassword,
				JWTSecret:    tt.jwtSecret,
			}
			server := &Server{
				Config: cfg,
				Logger: logger,
			}

			var body []byte
			if s, ok := tt.requestBody.(string); ok {
				body = []byte(s)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			server.Login(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectToken {
				var resp LoginResponse
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Token == "" {
					t.Error("expected token, got empty string")
				}
			}
		})
	}
}

func TestAuthCheck(t *testing.T) {
	cfg := &config.Config{
		AuthLogin: "admin",
	}
	server := &Server{
		Config: cfg,
	}

	req := httptest.NewRequest("GET", "/api/auth/check", nil)
	rr := httptest.NewRecorder()

	server.AuthCheck(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["authenticated"] != true {
		t.Errorf("expected authenticated true, got %v", resp["authenticated"])
	}
	if resp["user"] != "admin" {
		t.Errorf("expected user admin, got %v", resp["user"])
	}
}
