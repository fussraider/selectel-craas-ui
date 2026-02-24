package api

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/generic/selectel-craas-web/internal/craas"
	"github.com/go-chi/chi/v5/middleware"
)

// RequestLogger middleware logs request details.
func (s *Server) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		s.Logger.Debug("request started",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		next.ServeHTTP(ww, r)

		s.Logger.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"duration", time.Since(start),
		)
	})
}

// EnableCORS middleware sets CORS headers.
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ExecuteWithRetry executes an operation that requires a project token,
// handling token invalidation and retries automatically.
func (s *Server) ExecuteWithRetry(ctx context.Context, pid string, op func(token string) error) error {
	// 1. Get initial token
	token, err := s.Auth.GetProjectToken(pid)
	if err != nil {
		// If getting token fails, try invalidating and getting fresh one
		s.Logger.Warn("failed to get project token, retrying", "project_id", pid, "error", err)
		s.Auth.InvalidateProjectToken(pid)
		token, err = s.Auth.GetProjectToken(pid)
		if err != nil {
			return err
		}
	}

	// 2. Execute operation
	err = op(token)
	if err == nil {
		return nil
	}

	// 3. Check if the error is auth-related (401 Unauthorized)
	// We only retry if we suspect the token is invalid/expired.
	isAuthError := errors.Is(err, craas.ErrUnauthorized) || strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "Unauthorized")

	if isAuthError {
		s.Logger.Warn("auth error detected, retrying with token invalidation", "project_id", pid, "error", err)
		s.Auth.InvalidateProjectToken(pid)
		token, err = s.Auth.GetProjectToken(pid)
		if err != nil {
			return err // Failed to get fresh token
		}
		return op(token)
	}

	// For other errors, return immediately without retry
	return err
}
