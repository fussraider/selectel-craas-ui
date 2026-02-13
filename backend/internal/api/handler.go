package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/generic/selectel-craas-web/internal/auth"
	"github.com/generic/selectel-craas-web/internal/craas"
)

type Server struct {
	Auth   *auth.Client
	Craas  *craas.Service
	Logger *slog.Logger
}

func New(auth *auth.Client, craas *craas.Service, logger *slog.Logger) *chi.Mux {
	s := &Server{
		Auth:   auth,
		Craas:  craas,
		Logger: logger.With("service", "api"),
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(enableCORS)
	r.Use(s.RequestLogger) // Custom Logger

	r.Get("/api/auth/status", s.AuthStatus)
	r.Get("/api/projects", s.ListProjects)
	r.Get("/api/projects/{pid}/registries", s.ListRegistries)
	r.Delete("/api/projects/{pid}/registries/{rid}", s.DeleteRegistry)

	r.Get("/api/projects/{pid}/registries/{rid}/repositories", s.ListRepositories)
	r.Delete("/api/projects/{pid}/registries/{rid}/repository", s.DeleteRepository)

	r.Get("/api/projects/{pid}/registries/{rid}/images", s.ListImages)
	r.Delete("/api/projects/{pid}/registries/{rid}/images/{digest}", s.DeleteImage)
	r.Post("/api/projects/{pid}/registries/{rid}/cleanup", s.CleanupRepository)

	r.Get("/api/projects/{pid}/registries/{rid}/tags", s.ListTags)

	return r
}

func enableCORS(next http.Handler) http.Handler {
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

func (s *Server) AuthStatus(w http.ResponseWriter, r *http.Request) {
	_, err := s.Auth.GetAccountToken()
	if err != nil {
		s.Logger.Warn("auth check failed", "error", err)
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "authenticated"}`))
}

func (s *Server) ListProjects(w http.ResponseWriter, r *http.Request) {
	token, err := s.Auth.GetAccountToken()
	if err != nil {
		s.Logger.Error("failed to get account token", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projects, err := s.Auth.ListProjects(token)
	if err != nil {
		s.Logger.Warn("failed to list projects, retrying with token invalidation", "error", err)
		s.Auth.InvalidateAccountToken()
		token, err = s.Auth.GetAccountToken()
		if err == nil {
			projects, err = s.Auth.ListProjects(token)
		}
	}
	if err != nil {
		s.Logger.Error("failed to list projects after retry", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(projects)
}

func (s *Server) getProjectTokenWithRetry(pid string) (string, error) {
	token, err := s.Auth.GetProjectToken(pid)
	if err != nil {
		s.Logger.Warn("failed to get project token, retrying", "project_id", pid, "error", err)
		s.Auth.InvalidateProjectToken(pid)
		return s.Auth.GetProjectToken(pid)
	}
	return token, nil
}

func (s *Server) ListRegistries(w http.ResponseWriter, r *http.Request) {
	pid := chi.URLParam(r, "pid")
	s.Logger.Debug("listing registries request", "project_id", pid)

	token, err := s.getProjectTokenWithRetry(pid)
	if err != nil {
		s.Logger.Error("failed to get project token for registries list", "project_id", pid, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	registries, err := s.Craas.ListRegistries(r.Context(), token)
	if err != nil {
		s.Logger.Warn("failed to list registries, retrying", "project_id", pid, "error", err)
		s.Auth.InvalidateProjectToken(pid)
		token, err = s.Auth.GetProjectToken(pid)
		if err == nil {
			registries, err = s.Craas.ListRegistries(r.Context(), token)
		}
	}
	if err != nil {
		s.Logger.Error("failed to list registries after retry", "project_id", pid, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(registries)
}

func (s *Server) DeleteRegistry(w http.ResponseWriter, r *http.Request) {
	pid := chi.URLParam(r, "pid")
	rid := chi.URLParam(r, "rid")
	s.Logger.Info("deleting registry request", "project_id", pid, "registry_id", rid)

	token, err := s.getProjectTokenWithRetry(pid)
	if err != nil {
		s.Logger.Error("failed to get project token for registry deletion", "project_id", pid, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.Craas.DeleteRegistry(r.Context(), token, rid); err != nil {
		s.Logger.Warn("failed to delete registry, retrying", "registry_id", rid, "error", err)
		s.Auth.InvalidateProjectToken(pid)
		token, _ = s.Auth.GetProjectToken(pid)
		err = s.Craas.DeleteRegistry(r.Context(), token, rid)
	}

	if err != nil {
		s.Logger.Error("failed to delete registry after retry", "registry_id", rid, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) CleanupRepository(w http.ResponseWriter, r *http.Request) {
	pid := chi.URLParam(r, "pid")
	rid := chi.URLParam(r, "rid")
	rname := r.URL.Query().Get("repository")

	if rname == "" {
		http.Error(w, "repository param required", http.StatusBadRequest)
		return
	}

	var req craas.CleanupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Digests) == 0 && len(req.Tags) == 0 {
		http.Error(w, "digests or tags required", http.StatusBadRequest)
		return
	}

	token, err := s.getProjectTokenWithRetry(pid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := s.Craas.CleanupRepository(r.Context(), token, rid, rname, req.Digests, req.DisableGC)
	if err != nil {
		s.Logger.Warn("failed to cleanup repository, retrying", "error", err)
		s.Auth.InvalidateProjectToken(pid)
		token, _ = s.Auth.GetProjectToken(pid)
		result, err = s.Craas.CleanupRepository(r.Context(), token, rid, rname, req.Digests, req.DisableGC)
	}

	if err != nil {
		s.Logger.Error("failed to cleanup repository", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) ListRepositories(w http.ResponseWriter, r *http.Request) {
	pid := chi.URLParam(r, "pid")
	rid := chi.URLParam(r, "rid")

	token, err := s.getProjectTokenWithRetry(pid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	repos, err := s.Craas.ListRepositories(r.Context(), token, rid)
	if err != nil {
		s.Auth.InvalidateProjectToken(pid)
		token, _ = s.Auth.GetProjectToken(pid)
		repos, err = s.Craas.ListRepositories(r.Context(), token, rid)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(repos)
}

func (s *Server) DeleteRepository(w http.ResponseWriter, r *http.Request) {
	pid := chi.URLParam(r, "pid")
	rid := chi.URLParam(r, "rid")
	rname := r.URL.Query().Get("name")
	if rname == "" {
		http.Error(w, "repository name required", http.StatusBadRequest)
		return
	}

	token, err := s.getProjectTokenWithRetry(pid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.Craas.DeleteRepository(r.Context(), token, rid, rname); err != nil {
		s.Auth.InvalidateProjectToken(pid)
		token, _ = s.Auth.GetProjectToken(pid)
		err = s.Craas.DeleteRepository(r.Context(), token, rid, rname)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) ListImages(w http.ResponseWriter, r *http.Request) {
	pid := chi.URLParam(r, "pid")
	rid := chi.URLParam(r, "rid")
	rname := r.URL.Query().Get("repository")
	if rname == "" {
		http.Error(w, "repository param required", http.StatusBadRequest)
		return
	}

	token, err := s.getProjectTokenWithRetry(pid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	images, err := s.Craas.ListImages(r.Context(), token, rid, rname)
	if err != nil {
		s.Auth.InvalidateProjectToken(pid)
		token, _ = s.Auth.GetProjectToken(pid)
		images, err = s.Craas.ListImages(r.Context(), token, rid, rname)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(images)
}

func (s *Server) ListTags(w http.ResponseWriter, r *http.Request) {
	pid := chi.URLParam(r, "pid")
	rid := chi.URLParam(r, "rid")
	rname := r.URL.Query().Get("repository")
	if rname == "" {
		http.Error(w, "repository param required", http.StatusBadRequest)
		return
	}

	token, err := s.getProjectTokenWithRetry(pid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tags, err := s.Craas.ListTags(r.Context(), token, rid, rname)
	if err != nil {
		s.Auth.InvalidateProjectToken(pid)
		token, _ = s.Auth.GetProjectToken(pid)
		tags, err = s.Craas.ListTags(r.Context(), token, rid, rname)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tags)
}

func (s *Server) DeleteImage(w http.ResponseWriter, r *http.Request) {
	pid := chi.URLParam(r, "pid")
	rid := chi.URLParam(r, "rid")
	digest := chi.URLParam(r, "digest")
	rname := r.URL.Query().Get("repository")
	if rname == "" {
		http.Error(w, "repository param required", http.StatusBadRequest)
		return
	}

	token, err := s.getProjectTokenWithRetry(pid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.Craas.DeleteImage(r.Context(), token, rid, rname, digest); err != nil {
		s.Auth.InvalidateProjectToken(pid)
		token, _ = s.Auth.GetProjectToken(pid)
		err = s.Craas.DeleteImage(r.Context(), token, rid, rname, digest)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
