package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/mati/go-ticket/internal/api"
	"github.com/mati/go-ticket/internal/auth"
)

func AuthMiddleware(jwtService *auth.JWTService, next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			api.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			api.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")

		if bearerToken == "" {
			api.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		claims, err := jwtService.VerifyToken(bearerToken)
		if err != nil {
			slog.Warn("Authentication failed", "ip", r.RemoteAddr, "path", r.URL.Path, "method", r.Method, "error", err)
			api.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next(w, r.WithContext(ctx))
	}
}
