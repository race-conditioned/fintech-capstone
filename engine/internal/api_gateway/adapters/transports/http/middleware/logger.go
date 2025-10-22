package middleware

import (
	"fintech-capstone/m/v2/internal/api_gateway/adapters/transports/http/http_kit/writer"
	"fintech-capstone/m/v2/internal/api_gateway/ports"
	"net/http"
	"time"
)

// RequestLogger logs each HTTP request with method, path, status, bytes, latency, and client info.
func RequestLogger(logger ports.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			sw := writer.Wrap(w)

			next.ServeHTTP(sw, r)

			latency := time.Since(start)
			status, bytes, _ := writer.Status(sw)
			rid, _ := FromContext(r.Context())

			logger.Info("http request",
				ports.Field{Key: "method", Value: r.Method},
				ports.Field{Key: "path", Value: r.URL.Path},
				ports.Field{Key: "status", Value: status},
				ports.Field{Key: "bytes", Value: bytes},
				ports.Field{Key: "latency_ms", Value: latency.Milliseconds()},
				ports.Field{Key: "client_id", Value: r.Header.Get("X-Client-ID")},
				ports.Field{Key: "remote_addr", Value: r.RemoteAddr},
				ports.Field{Key: "user_agent", Value: r.UserAgent()},
				ports.Field{Key: "request_id", Value: rid},
			)
		})
	}
}
