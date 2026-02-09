package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/mati/go-ticket/internal/api"
	"github.com/mati/go-ticket/internal/auth"
	"github.com/mati/go-ticket/internal/domain"
)

type userData struct {
	ID   uuid.UUID
	Role domain.UserRole
}

type contextKey string

const userContextKey contextKey = "user"

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

		var user = userData{
			ID:   claims.UserID,
			Role: claims.Role,
		}
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next(w, r.WithContext(ctx))
	}
}

func RequireRole(allowedRoles []domain.UserRole, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user, ok := r.Context().Value(userContextKey).(userData)
		if !ok {
			api.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if slices.Contains(allowedRoles, user.Role) {
			next(w, r)
			return
		}
		api.ResponseError(w, http.StatusForbidden, "Forbidden")
	}
}
