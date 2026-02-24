package middleware

import (
	"log/slog"
	"net/http"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error( //nolint:gosec // G706: slog does not concatenate strings
					"Recovered from panic",
					"panic", rec,
					"path", r.URL.Path,
					"method", r.Method,
				)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
