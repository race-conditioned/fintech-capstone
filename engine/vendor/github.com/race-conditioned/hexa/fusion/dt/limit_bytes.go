package dt

import (
	"net/http"

	"github.com/race-conditioned/hexa/fusion/intake"
)

// limitBytes limits the size of request bodies to n bytes.
func limitBytes(next http.Handler, n int64, hookFn intake.Hook) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limited := http.MaxBytesReader(w, r.Body, n)
		r.Body = limited

		if hookFn != nil {
			hookFn(intake.Event{
				Protocol: "http",
				Target:   r.URL.Path,
				ClientID: extractClientID(r),
			})
		}

		next.ServeHTTP(w, r)
	})
}
