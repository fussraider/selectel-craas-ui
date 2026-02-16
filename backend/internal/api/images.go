package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) ListImages(w http.ResponseWriter, r *http.Request) {
	pid := chi.URLParam(r, "pid")
	rid := chi.URLParam(r, "rid")
	rname := r.URL.Query().Get("repository")
	if rname == "" {
		http.Error(w, "repository param required", http.StatusBadRequest)
		return
	}

	var result interface{}
	err := s.ExecuteWithRetry(r.Context(), pid, func(token string) error {
		var err error
		result, err = s.Craas.ListImages(r.Context(), token, rid, rname)
		return err
	})

	if err != nil {
		s.Logger.Error("failed to list images", "registry_id", rid, "repository", rname, "error", err)
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	RespondJSON(w, http.StatusOK, result)
}

func (s *Server) ListTags(w http.ResponseWriter, r *http.Request) {
	pid := chi.URLParam(r, "pid")
	rid := chi.URLParam(r, "rid")
	rname := r.URL.Query().Get("repository")
	if rname == "" {
		http.Error(w, "repository param required", http.StatusBadRequest)
		return
	}

	var result interface{}
	err := s.ExecuteWithRetry(r.Context(), pid, func(token string) error {
		var err error
		result, err = s.Craas.ListTags(r.Context(), token, rid, rname)
		return err
	})

	if err != nil {
		s.Logger.Error("failed to list tags", "registry_id", rid, "repository", rname, "error", err)
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	RespondJSON(w, http.StatusOK, result)
}

func (s *Server) DeleteImage(w http.ResponseWriter, r *http.Request) {
	if !s.checkDeleteImage(w) {
		return
	}

	pid := chi.URLParam(r, "pid")
	rid := chi.URLParam(r, "rid")
	digest := chi.URLParam(r, "digest")
	rname := r.URL.Query().Get("repository")
	if rname == "" {
		http.Error(w, "repository param required", http.StatusBadRequest)
		return
	}

	err := s.ExecuteWithRetry(r.Context(), pid, func(token string) error {
		return s.Craas.DeleteImage(r.Context(), token, rid, rname, digest)
	})

	if err != nil {
		s.Logger.Error("failed to delete image", "registry_id", rid, "repository", rname, "digest", digest, "error", err)
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
