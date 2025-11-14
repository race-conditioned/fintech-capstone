package dt

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/race-conditioned/hexa/fusion/intake"
)

type ctxKey int

const requestIDKey ctxKey = iota

// requestID is a middleware that ensures each request has a unique request ID.
// If the incoming request has an "X-Request-ID" header, it uses that value.
// Otherwise, it generates a new random ID.
func requestID(next http.Handler, hookFn intake.RequestIDHook) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			var buf [16]byte
			_, _ = rand.Read(buf[:])
			id = hex.EncodeToString(buf[:])
		}
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		w.Header().Set("X-Request-ID", id)
		if hookFn != nil {
			hookFn(intake.Event{
				Protocol: "http",
				Target:   r.URL.Path,
				ClientID: extractClientID(r),
			}, id)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromContext retrieves the request ID from the context.
func FromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey).(string)
	return id, ok
}
