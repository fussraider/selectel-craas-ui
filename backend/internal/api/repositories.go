package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/generic/selectel-craas-web/internal/craas"
)

func (s *Server) ListRepositories(w http.ResponseWriter, r *http.Request) {
	pid := chi.URLParam(r, "pid")
	rid := chi.URLParam(r, "rid")

	var result interface{}
	err := s.ExecuteWithRetry(r.Context(), pid, func(token string) error {
		var err error
		result, err = s.Craas.ListRepositories(r.Context(), token, rid)
		return err
	})

	if err != nil {
		s.Logger.Error("failed to list repositories", "registry_id", rid, "error", err)
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	RespondJSON(w, http.StatusOK, result)
}

func (s *Server) DeleteRepository(w http.ResponseWriter, r *http.Request) {
	if !s.checkDeleteRepository(w) {
		return
	}

	pid := chi.URLParam(r, "pid")
	rid := chi.URLParam(r, "rid")
	rname := r.URL.Query().Get("name")
	if rname == "" {
		http.Error(w, "repository name required", http.StatusBadRequest)
		return
	}

	err := s.ExecuteWithRetry(r.Context(), pid, func(token string) error {
		return s.Craas.DeleteRepository(r.Context(), token, rid, rname)
	})

	if err != nil {
		s.Logger.Error("failed to delete repository", "registry_id", rid, "repository", rname, "error", err)
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) CleanupRepository(w http.ResponseWriter, r *http.Request) {
	if !s.checkDeleteImage(w) {
		return
	}

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

	// Safety check: Protected Tags
	if len(s.Config.ProtectedTags) > 0 && len(req.Digests) > 0 {
		// We need to fetch images to check their tags against protected tags
		// because the request only contains digests.
		var protectedFound bool
		var protectedTag string

		err := s.ExecuteWithRetry(r.Context(), pid, func(token string) error {
			images, err := s.Craas.ListImages(r.Context(), token, rid, rname)
			if err != nil {
				return err
			}

			// Build map for quick lookup
			digestTags := make(map[string][]string)
			for _, img := range images {
				digestTags[img.Digest] = img.Tags
			}

			// Check requested digests
			for _, digest := range req.Digests {
				if tags, ok := digestTags[digest]; ok {
					for _, tag := range tags {
						for _, protected := range s.Config.ProtectedTags {
							if tag == protected {
								protectedFound = true
								protectedTag = tag
								return nil // Stop checking
							}
						}
					}
				}
			}
			return nil
		})

		if err != nil {
			s.Logger.Error("failed to verify protected tags", "error", err)
			RespondError(w, http.StatusInternalServerError, err)
			return
		}

		if protectedFound {
			http.Error(w, fmt.Sprintf("cannot delete image with protected tag: %s", protectedTag), http.StatusForbidden)
			return
		}
	}

	var result interface{}
	err := s.ExecuteWithRetry(r.Context(), pid, func(token string) error {
		var err error
		result, err = s.Craas.CleanupRepository(r.Context(), token, rid, rname, req.Digests, req.DisableGC)
		return err
	})

	if err != nil {
		s.Logger.Error("failed to cleanup repository", "registry_id", rid, "repository", rname, "error", err)
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	RespondJSON(w, http.StatusOK, result)
}
