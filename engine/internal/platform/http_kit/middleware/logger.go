package middleware

import (
	"fintech-capstone/m/v2/internal/platform"
	"fintech-capstone/m/v2/internal/platform/http_kit/writer"
	"net/http"
	"time"
)

// RequestLogger logs each HTTP request with method, path, status, bytes, latency, and client info.
func RequestLogger(logger platform.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			sw := writer.Wrap(w)

			next.ServeHTTP(sw, r)

			latency := time.Since(start)
			status, bytes, _ := writer.Status(sw)
			rid, _ := FromContext(r.Context())

			logger.Info("http request",
				platform.Field{Key: "method", Value: r.Method},
				platform.Field{Key: "path", Value: r.URL.Path},
				platform.Field{Key: "status", Value: status},
				platform.Field{Key: "bytes", Value: bytes},
				platform.Field{Key: "latency_ms", Value: latency.Milliseconds()},
				platform.Field{Key: "client_id", Value: r.Header.Get("X-Client-ID")},
				platform.Field{Key: "remote_addr", Value: r.RemoteAddr},
				platform.Field{Key: "user_agent", Value: r.UserAgent()},
				platform.Field{Key: "request_id", Value: rid},
			)
		})
	}
}
