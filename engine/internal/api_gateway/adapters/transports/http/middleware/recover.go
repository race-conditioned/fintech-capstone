package middleware

import (
	"fintech-capstone/m/v2/internal/api_gateway/adapters/transports/http/http_kit/writer"
	"fintech-capstone/m/v2/internal/api_gateway/ports"
	"fmt"
	"net/http"
	"runtime/debug"
)

// RecovererWithLogger recovers from panics, logs details, and returns 500.
func RecovererWithLogger(logger ports.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					stack := debug.Stack()

					logger.Error(fmt.Errorf("panic: %v", rec),
						ports.Field{Key: "stack", Value: string(stack)},
						ports.Field{Key: "method", Value: r.Method},
						ports.Field{Key: "path", Value: r.URL.Path},
						ports.Field{Key: "remote_addr", Value: r.RemoteAddr},
					)

					writer.Error(w, http.StatusInternalServerError, "internal server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
