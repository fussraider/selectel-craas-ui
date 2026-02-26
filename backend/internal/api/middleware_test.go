package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/generic/selectel-craas-web/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

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
