package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		record := &ResponseRecord{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}

		next.ServeHTTP(record, r)

		duration := time.Since(start)

		slog.Info("request completed", //nolint:gosec // G706: slog uses structured fields
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path), //nolint:gosec // G706: slog uses structured fields
			slog.Duration("duration", duration),
			slog.Int("status", record.Status),
		)
	})
}

type ResponseRecord struct {
	http.ResponseWriter
	Status  int
	Written bool
}

func (r *ResponseRecord) WriteHeader(status int) {
	if r.Written {
		return
	}
	r.Status = status
	r.Written = true
	r.ResponseWriter.WriteHeader(status)
}

func (r *ResponseRecord) Write(b []byte) (int, error) {
	if !r.Written {
		r.WriteHeader(http.StatusOK)
	}
	return r.ResponseWriter.Write(b)
}
