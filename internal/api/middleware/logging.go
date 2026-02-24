package middleware

import (
	"log/slog"
	"net/http"
	"strings"
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

		safePath := strings.ReplaceAll(r.URL.Path, "\n", "\\n")
		safePath = strings.ReplaceAll(safePath, "\r", "\\r")

		slog.Info("request completed",
			slog.String("method", r.Method),
			slog.String("path", safePath),
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
