package middleware

import (
	"log/slog"
	"net/http"

	"github.com/mati/go-ticket/internal/ratelimit"
)

func RateLimiterMiddleware(
	limiter *ratelimit.RateLimiter,
	getKey func(r *http.Request) string,
) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			key := getKey(r)
			allowed, err := limiter.Allow(r.Context(), key)
			if err != nil {
				slog.Error("Rate limiter error", "error", err)
				next(w, r)
				return
			}
			if !allowed {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			next(w, r)
		}
	}
}
