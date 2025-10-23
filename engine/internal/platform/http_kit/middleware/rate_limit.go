package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"
)

// AllowFunc is a minimal adapter so you can plug in any limiter.
type AllowFunc func(ctx context.Context, key string) bool

// ClientKeyFromRequest picks API key, then X-Forwarded-For, then RemoteAddr.
func ClientKeyFromRequest(r *http.Request) string {
	if k := r.Header.Get("X-API-Key"); k != "" {
		return k
	}
	if k := r.Header.Get("X-Client-ID"); k != "" {
		return k
	}
	ip := clientIP(r)
	if ip != "" {
		return ip
	}
	return r.RemoteAddr
}

func clientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// RateLimitHTTP rejects early (429) without reading the body.
func RateLimitHTTP(allow AllowFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Very short budget so we never block at this layer.
			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Millisecond)
			defer cancel()

			key := ClientKeyFromRequest(r)
			if !allow(ctx, key) {
				// Best-effort: close body to free the conn; we didn't read it.
				_ = r.Body.Close()
				w.Header().Set("Retry-After", "1")
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
