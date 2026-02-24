package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/generic/selectel-craas-web/internal/config"
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
