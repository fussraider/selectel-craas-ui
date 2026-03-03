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
		cookieSecure   bool
		cookieSameSite string
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
			name:           "Invalid Password",
			authEnabled:    true,
			authLogin:      "admin",
			authPassword:   "password",
			requestBody:    LoginRequest{Login: "admin", Password: "wrong-password"},
			expectedStatus: http.StatusUnauthorized,
			expectCookie:   false,
		},
		{
			name:           "Invalid Login",
			authEnabled:    true,
			authLogin:      "admin",
			authPassword:   "password",
			requestBody:    LoginRequest{Login: "wrong-user", Password: "password"},
			expectedStatus: http.StatusUnauthorized,
			expectCookie:    false,
		},
		{
			name:           "Successful Login",
			authEnabled:    true,
			authLogin:      "admin",
			authPassword:   "password",
			jwtSecret:      "secret",
			cookieSecure:   true,
			cookieSameSite: "lax",
			requestBody:    LoginRequest{Login: "admin", Password: "password"},
			expectedStatus: http.StatusOK,
			expectCookie:   true,
		},
		{
			name:           "Successful Login Insecure Cookie",
			authEnabled:    true,
			authLogin:      "admin",
			authPassword:   "password",
			jwtSecret:      "secret",
			cookieSecure:   false,
			cookieSameSite: "none",
			requestBody:    LoginRequest{Login: "admin", Password: "password"},
			expectedStatus: http.StatusOK,
			expectCookie:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				AuthEnabled:    tt.authEnabled,
				AuthLogin:      tt.authLogin,
				AuthPassword:   tt.authPassword,
				JWTSecret:      tt.jwtSecret,
				CookieSecure:   tt.cookieSecure,
				CookieSameSite: tt.cookieSameSite,
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
						if cookie.Secure != tt.cookieSecure {
							t.Errorf("expected Secure cookie: %v, got %v", tt.cookieSecure, cookie.Secure)
						}
						expectedSameSite := http.SameSiteLaxMode
						if tt.cookieSameSite == "none" {
							expectedSameSite = http.SameSiteNoneMode
						} else if tt.cookieSameSite == "strict" {
							expectedSameSite = http.SameSiteStrictMode
						}
						if cookie.SameSite != expectedSameSite {
							t.Errorf("expected SameSite cookie: %v, got %v", expectedSameSite, cookie.SameSite)
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
	server := &Server{
		Config: &config.Config{
			CookieSecure: true,
		},
	}
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
