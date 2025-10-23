package middleware

import (
	"net/http"
)

// MaxInFlight limits the number of concurrent in-flight HTTP requests to n.
func MaxInFlight(n int) func(http.Handler) http.Handler {
	sem := make(chan struct{}, n)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
				next.ServeHTTP(w, r)
			default:
				http.Error(w, "server busy", http.StatusServiceUnavailable)
			}
		})
	}
}
