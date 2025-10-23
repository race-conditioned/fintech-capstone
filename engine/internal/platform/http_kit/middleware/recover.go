package middleware

import (
	"fintech-capstone/m/v2/internal/platform"
	"fintech-capstone/m/v2/internal/platform/http_kit/writer"
	"fmt"
	"net/http"
	"runtime/debug"
)

// RecovererWithLogger recovers from panics, logs details, and returns 500.
func RecovererWithLogger(logger platform.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					stack := debug.Stack()

					logger.Error(fmt.Errorf("panic: %v", rec),
						platform.Field{Key: "stack", Value: string(stack)},
						platform.Field{Key: "method", Value: r.Method},
						platform.Field{Key: "path", Value: r.URL.Path},
						platform.Field{Key: "remote_addr", Value: r.RemoteAddr},
					)

					writer.Error(w, http.StatusInternalServerError, "internal server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
