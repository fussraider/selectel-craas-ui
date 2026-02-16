package api

import (
	"log/slog"

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
	r.Use(EnableCORS)
	r.Use(s.RequestLogger)

	// Projects
	r.Get("/api/auth/status", s.AuthStatus)
	r.Get("/api/projects", s.ListProjects)

	// Registries
	r.Get("/api/projects/{pid}/registries", s.ListRegistries)
	r.Delete("/api/projects/{pid}/registries/{rid}", s.DeleteRegistry)
	r.Get("/api/projects/{pid}/registries/{rid}/gc", s.GetGCInfo)
	r.Post("/api/projects/{pid}/registries/{rid}/gc", s.StartGC)

	// Repositories
	r.Get("/api/projects/{pid}/registries/{rid}/repositories", s.ListRepositories)
	r.Delete("/api/projects/{pid}/registries/{rid}/repository", s.DeleteRepository)
	r.Post("/api/projects/{pid}/registries/{rid}/cleanup", s.CleanupRepository)

	// Images
	r.Get("/api/projects/{pid}/registries/{rid}/images", s.ListImages)
	r.Delete("/api/projects/{pid}/registries/{rid}/images/{digest}", s.DeleteImage)
	r.Get("/api/projects/{pid}/registries/{rid}/tags", s.ListTags)

	return r
}
