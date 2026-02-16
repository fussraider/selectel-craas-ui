package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) ListRegistries(w http.ResponseWriter, r *http.Request) {
	pid := chi.URLParam(r, "pid")
	s.Logger.Debug("listing registries request", "project_id", pid)

	var result interface{}
	err := s.ExecuteWithRetry(r.Context(), pid, func(token string) error {
		var err error
		result, err = s.Craas.ListRegistries(r.Context(), token)
		return err
	})

	if err != nil {
		s.Logger.Error("failed to list registries", "project_id", pid, "error", err)
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	RespondJSON(w, http.StatusOK, result)
}

func (s *Server) DeleteRegistry(w http.ResponseWriter, r *http.Request) {
	if !s.checkDeleteRegistry(w) {
		return
	}

	pid := chi.URLParam(r, "pid")
	rid := chi.URLParam(r, "rid")
	s.Logger.Info("deleting registry request", "project_id", pid, "registry_id", rid)

	err := s.ExecuteWithRetry(r.Context(), pid, func(token string) error {
		return s.Craas.DeleteRegistry(r.Context(), token, rid)
	})

	if err != nil {
		s.Logger.Error("failed to delete registry", "registry_id", rid, "error", err)
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) GetGCInfo(w http.ResponseWriter, r *http.Request) {
	pid := chi.URLParam(r, "pid")
	rid := chi.URLParam(r, "rid")

	var result interface{}
	err := s.ExecuteWithRetry(r.Context(), pid, func(token string) error {
		var err error
		result, err = s.Craas.GetGCInfo(r.Context(), token, rid)
		return err
	})

	if err != nil {
		s.Logger.Error("failed to get GC info", "registry_id", rid, "error", err)
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	RespondJSON(w, http.StatusOK, result)
}

func (s *Server) StartGC(w http.ResponseWriter, r *http.Request) {
	pid := chi.URLParam(r, "pid")
	rid := chi.URLParam(r, "rid")

	err := s.ExecuteWithRetry(r.Context(), pid, func(token string) error {
		return s.Craas.StartGC(r.Context(), token, rid)
	})

	if err != nil {
		if err.Error() == "garbage collection already in progress" {
			RespondError(w, http.StatusConflict, err)
			return
		}
		s.Logger.Error("failed to start GC", "registry_id", rid, "error", err)
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
