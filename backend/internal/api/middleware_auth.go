package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var ErrUnauthorized = errors.New("unauthorized")

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.AuthEnabled {
			next.ServeHTTP(w, r)
			return
		}

		var tokenString string

		// First try to get token from cookie
		if cookie, err := r.Cookie("auth_token"); err == nil {
			tokenString = cookie.Value
		}

		// Fallback to Authorization header
		if tokenString == "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}
		}

		if tokenString == "" {
			RespondError(w, http.StatusUnauthorized, ErrUnauthorized)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(s.Config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			RespondError(w, http.StatusUnauthorized, ErrUnauthorized)
			return
		}

		// Optionally extract claims and put in context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx := context.WithValue(r.Context(), "user", claims["sub"])
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			RespondError(w, http.StatusUnauthorized, ErrUnauthorized)
		}
	})
}
