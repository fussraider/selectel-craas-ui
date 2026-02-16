package api

import (
	"net/http"
)

func (s *Server) AuthStatus(w http.ResponseWriter, r *http.Request) {
	_, err := s.Auth.GetAccountToken()
	if err != nil {
		s.Logger.Warn("auth check failed", "error", err)
		RespondError(w, http.StatusUnauthorized, err)
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"status": "authenticated"})
}

func (s *Server) ListProjects(w http.ResponseWriter, r *http.Request) {
	token, err := s.Auth.GetAccountToken()
	if err != nil {
		s.Logger.Error("failed to get account token", "error", err)
		RespondError(w, http.StatusInternalServerError, err)
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
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	RespondJSON(w, http.StatusOK, projects)
}
