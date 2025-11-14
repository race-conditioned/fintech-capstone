package dt

import (
	"context"
	"net/http"
	"time"

	"github.com/race-conditioned/hexa/fusion/intake"
)

// rateLimit rejects early (429) without reading the body.
// It is purely transport-level (does not depend on business logic).
func rateLimit(next http.Handler, hook intake.Hook, allow intake.AllowFunc) http.Handler {
	// No limiter configured → no-op
	if allow == nil {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Short timeout budget: don’t block the transport thread.
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Millisecond)
		defer cancel()

		key := clientKeyFromRequest(r)

		if !allow(ctx, key) {
			_ = r.Body.Close()
			w.Header().Set("Retry-After", "1")
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)

			if hook != nil {
				hook(intake.Event{
					Protocol: "http",
					Target:   r.URL.Path,
					ClientID: key,
				})
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

// clientKeyFromRequest extracts a usable client identifier.
func clientKeyFromRequest(r *http.Request) string {
	if k := r.Header.Get("X-API-Key"); k != "" {
		return k
	}
	if k := r.Header.Get("X-Client-ID"); k != "" {
		return k
	}
	return r.RemoteAddr
}
