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

		slog.Info("request completed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Duration("duration", duration),
			slog.Int("status", record.Status),
		)
	})
}

type ResponseRecord struct {
	http.ResponseWriter
	Status  int
	Writter bool
}

func (r *ResponseRecord) WriteHeader(status int) {
	if r.Writter {
		return
	}
	r.Status = status
	r.Writter = true
	r.ResponseWriter.WriteHeader(status)
}

func (r *ResponseRecord) Write(b []byte) (int, error) {
	if !r.Writter {
		r.WriteHeader(http.StatusOK)
	}
	return r.ResponseWriter.Write(b)
}
