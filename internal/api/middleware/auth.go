package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/mati/go-ticket/internal/auth"
	"github.com/mati/go-ticket/internal/domain"
)

type userData struct {
	ID    uuid.UUID
	Role  domain.UserRole
	Email string
}

type contextKey string

const userContextKey contextKey = "user"

func AuthMiddleware(jwtService *auth.JWTService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")

		if bearerToken == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := jwtService.VerifyToken(bearerToken)
		if err != nil {
			slog.Warn( //nolint:gosec // G706: slog uses structured fields
				"Authentication failed",
				"ip", r.RemoteAddr,
				"path", r.URL.Path,
				"method", r.Method,
				"error", err,
			)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var user = userData{
			ID:    claims.UserID,
			Role:  claims.Role,
			Email: claims.Email,
		}
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next(w, r.WithContext(ctx))
	}
}

func RequireRole(allowedRoles []domain.UserRole, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(userContextKey).(userData)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if slices.Contains(allowedRoles, user.Role) {
			next(w, r)
			return
		}
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func GetUserDataFromContext(ctx context.Context) (userData, bool) {
	user, ok := ctx.Value(userContextKey).(userData)
	if !ok {
		return userData{}, false
	}
	return user, true
}

func WithTestUser(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, userContextKey, userData{
		Email: email,
		ID:    uuid.New(),
	})
}
