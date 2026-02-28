package api

import (
	"crypto/subtle"
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
	User string `json:"user"`
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

	loginMatch := subtle.ConstantTimeCompare([]byte(req.Login), []byte(s.Config.AuthLogin))
	passMatch := subtle.ConstantTimeCompare([]byte(req.Password), []byte(s.Config.AuthPassword))

	if loginMatch&passMatch != 1 {
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

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	RespondJSON(w, http.StatusOK, LoginResponse{User: req.Login})
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusOK)
}

func (s *Server) AuthCheck(w http.ResponseWriter, r *http.Request) {
	// If the request reached here, it passed the AuthMiddleware (if enabled).
	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"authenticated": true,
		"user":          s.Config.AuthLogin,
	})
}
