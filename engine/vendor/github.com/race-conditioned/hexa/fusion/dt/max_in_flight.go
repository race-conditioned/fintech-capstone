package dt

import (
	"net/http"

	"github.com/race-conditioned/hexa/fusion/intake"
)

// maxInFlight limits the number of concurrent in-flight HTTP requests to n.
func maxInFlight(next http.Handler, n int, hookFn intake.Hook) http.Handler {
	sem := make(chan struct{}, n)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hookFn != nil {
			hookFn(intake.Event{
				Protocol: "http",
				Target:   r.URL.Path,
				ClientID: extractClientID(r),
			})
		}

		select {
		case sem <- struct{}{}:
			defer func() { <-sem }()
			next.ServeHTTP(w, r)
		default:
			http.Error(w, "server busy", http.StatusServiceUnavailable)
		}
	})
}
