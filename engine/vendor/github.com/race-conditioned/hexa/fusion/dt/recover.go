package dt

import (
	"net/http"

	"github.com/race-conditioned/hexa/fusion/dt/nolan"
	"github.com/race-conditioned/hexa/fusion/intake"
)

// recoverer recovers from panics, logs details, and returns 500.
func recoverer(next http.Handler, hookFn intake.RecoverHook) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				writer := nolan.NewWriter()
				if hookFn != nil {
					hookFn(intake.Event{
						Protocol: "http",
						Target:   r.URL.Path,
						ClientID: extractClientID(r),
					}, "") // TODO: why another arg?
				}

				writer.Error(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
