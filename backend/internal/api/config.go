package api

import (
	"fmt"
	"net/http"
)

func (s *Server) GetConfig(w http.ResponseWriter, r *http.Request) {
	// Expose only safe configuration
	cfg := map[string]bool{
		"enableDeleteRegistry":   s.Config.EnableDeleteRegistry,
		"enableDeleteRepository": s.Config.EnableDeleteRepository,
		"enableDeleteImage":      s.Config.EnableDeleteImage,
	}
	RespondJSON(w, http.StatusOK, cfg)
}

// Error forbidden
var ErrForbidden = fmt.Errorf("action disabled by configuration")

func (s *Server) checkDeleteRegistry(w http.ResponseWriter) bool {
	if !s.Config.EnableDeleteRegistry {
		RespondError(w, http.StatusForbidden, ErrForbidden)
		return false
	}
	return true
}

func (s *Server) checkDeleteRepository(w http.ResponseWriter) bool {
	if !s.Config.EnableDeleteRepository {
		RespondError(w, http.StatusForbidden, ErrForbidden)
		return false
	}
	return true
}

func (s *Server) checkDeleteImage(w http.ResponseWriter) bool {
	if !s.Config.EnableDeleteImage {
		RespondError(w, http.StatusForbidden, ErrForbidden)
		return false
	}
	return true
}
