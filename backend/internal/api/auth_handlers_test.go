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
		expectCookie   bool
	}{
		{
			name:           "Auth Disabled",
			authEnabled:    false,
			requestBody:    LoginRequest{Login: "user", Password: "password"},
			expectedStatus: http.StatusBadRequest,
			expectCookie:   false,
		},
		{
			name:           "Invalid JSON",
			authEnabled:    true,
			requestBody:    "invalid-json",
			expectedStatus: http.StatusBadRequest,
			expectCookie:   false,
		},
		{
			name:           "Invalid Credentials",
			authEnabled:    true,
			authLogin:      "admin",
			authPassword:   "password",
			requestBody:    LoginRequest{Login: "admin", Password: "wrong-password"},
			expectedStatus: http.StatusUnauthorized,
			expectCookie:   false,
		},
		{
			name:           "Successful Login",
			authEnabled:    true,
			authLogin:      "admin",
			authPassword:   "password",
			jwtSecret:      "secret",
			requestBody:    LoginRequest{Login: "admin", Password: "password"},
			expectedStatus: http.StatusOK,
			expectCookie:   true,
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

			if tt.expectCookie {
				var resp LoginResponse
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.User == "" {
					t.Error("expected user, got empty string")
				}

				cookies := rr.Result().Cookies()
				found := false
				for _, cookie := range cookies {
					if cookie.Name == "auth_token" {
						found = true
						if cookie.Value == "" {
							t.Error("expected auth_token value, got empty string")
						}
						if !cookie.HttpOnly {
							t.Error("expected HttpOnly cookie")
						}
						if !cookie.Secure {
							t.Error("expected Secure cookie")
						}
						break
					}
				}
				if !found {
					t.Error("expected auth_token cookie not found")
				}
			}
		})
	}
}

func TestLogout(t *testing.T) {
	server := &Server{}
	req := httptest.NewRequest("POST", "/api/logout", nil)
	rr := httptest.NewRecorder()

	server.Logout(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	cookies := rr.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			found = true
			if cookie.MaxAge != -1 {
				t.Errorf("expected MaxAge -1, got %d", cookie.MaxAge)
			}
			break
		}
	}
	if !found {
		t.Error("expected auth_token cookie to be cleared")
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
