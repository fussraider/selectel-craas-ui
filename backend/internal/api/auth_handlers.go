package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	if !s.Config.AuthEnabled {
		RespondError(w, http.StatusBadRequest, errors.New("authentication is disabled"))
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return
	}

	if req.Login != s.Config.AuthLogin || req.Password != s.Config.AuthPassword {
		RespondError(w, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": req.Login,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.Config.JWTSecret))
	if err != nil {
		s.Logger.Error("failed to sign token", "error", err)
		RespondError(w, http.StatusInternalServerError, err)
		return
	}

	RespondJSON(w, http.StatusOK, LoginResponse{Token: tokenString})
}

func (s *Server) AuthCheck(w http.ResponseWriter, r *http.Request) {
	// If the request reached here, it passed the AuthMiddleware (if enabled).
	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"authenticated": true,
		"user":          s.Config.AuthLogin,
	})
}
