package dt

import (
	"net/http"
	"time"

	"github.com/race-conditioned/hexa/fusion/dt/nolan"
	"github.com/race-conditioned/hexa/fusion/intake"
)

// requestLogger builds an HTTP middleware that logs request lifecycle info
// using the transport-neutral intake.LogHook from the ingress Spec.
func requestLogger(next http.Handler, hookFn intake.LogHook) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the ResponseWriter to capture status and bytes written.
		writer := nolan.NewWriter()
		sw := writer.Wrap(w)

		next.ServeHTTP(sw, r)

		// Collect neutral metadata for the hook.
		latency := time.Since(start)
		status, bytesWritten, _ := writer.Status(sw)
		rid, _ := FromContext(r.Context())

		if hookFn != nil {
			hookFn(
				intake.Event{
					Protocol: "http",
					Target:   r.URL.Path,
					ClientID: extractClientID(r),
				},
				intake.LogFields{
					StatusCode: status,
					Bytes:      int64(bytesWritten),
					LatencyMs:  latency.Milliseconds(),
					Method:     r.Method,
					RequestID:  rid,
					RemoteAddr: r.RemoteAddr,
					UserAgent:  r.UserAgent(),
				},
			)
		}
	})
}
