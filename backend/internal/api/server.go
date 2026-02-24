package api

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/generic/selectel-craas-web/internal/auth"
	"github.com/generic/selectel-craas-web/internal/config"
	"github.com/generic/selectel-craas-web/internal/craas"
)

type Server struct {
	Auth   *auth.Client
	Craas  *craas.Service
	Logger *slog.Logger
	Config *config.Config
}

func New(auth *auth.Client, craas *craas.Service, logger *slog.Logger, cfg *config.Config) *chi.Mux {
	s := &Server{
		Auth:   auth,
		Craas:  craas,
		Logger: logger.With("service", "api"),
		Config: cfg,
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(s.SecurityHeaders)
	r.Use(s.EnableCORS)
	r.Use(s.RequestLogger)

	// Public routes
	r.Get("/api/config", s.GetConfig)
	r.Post("/api/login", s.Login)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(s.AuthMiddleware)

		r.Get("/api/auth/check", s.AuthCheck)

		// Projects
		r.Get("/api/auth/status", s.AuthStatus) // Checks upstream auth status
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
	})

	return r
}
