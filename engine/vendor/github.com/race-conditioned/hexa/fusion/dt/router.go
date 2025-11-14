package dt

import (
	"fmt"
	"net/http"

	"github.com/race-conditioned/hexa/fusion/intake"
	"github.com/race-conditioned/hexa/horizon"
	"github.com/race-conditioned/hexa/horizon/ports/inbound"
)

type DTFusion[c inbound.Ctx] struct {
	gw     *horizon.Gateway[c]
	mux    *http.ServeMux
	spec   intake.Spec
	ctx    c
	routes []Route[c]
}

func NewFusion[c inbound.Ctx](
	ctx c,
	spec intake.Spec,
	gw *horizon.Gateway[c],
	routes []Route[c],
) DTFusion[c] {
	return DTFusion[c]{gw: gw, mux: http.NewServeMux(), spec: spec, ctx: ctx, routes: routes}
}

func (dt DTFusion[c]) Build() http.Handler {
	for _, rt := range dt.routes {
		h, ok := dt.gw.Handler(rt.HandlerKey)
		if !ok {
			continue
		}
		fmt.Println("found handler for: ", rt.HandlerKey.String(), h)

		path := rt.Path
		if path == "" {
			path = "/" + rt.HandlerKey.String()
		}
		fmt.Println("binding path: ", path)

		dt.mux.HandleFunc(path,
			Unary(
				h,
				rt.NewPayload,
				DefaultMeta,
				dt.ctx,
			),
		)
	}
	return dt.applySpec(dt.mux)
}

func (dt DTFusion[c]) applySpec(mux http.Handler) http.Handler {
	// TODO: allow the user to define order and mix in custom middleware
	if dt.spec.MaxBodyBytes > 0 {
		mux = limitBytes(mux, dt.spec.MaxBodyBytes, dt.spec.OnBodyTooLarge)
	}
	if dt.spec.MaxInFlight > 0 {
		mux = maxInFlight(mux, dt.spec.MaxInFlight, dt.spec.OnTooManyInFlight)
	}
	if dt.spec.EnableRecover {
		mux = recoverer(mux, dt.spec.OnRecover)
	}

	if dt.spec.EnableReqID {
		mux = requestID(mux, dt.spec.OnRequestID)
	}
	if dt.spec.EnableReqLog {
		mux = requestLogger(mux, dt.spec.OnLog)
	}
	if dt.spec.EnableRateLimiting {
		mux = rateLimit(mux, dt.spec.OnRateLimit, dt.spec.RateLimiter)
	}

	return chain(mux, typeMiddleware(dt.spec.Middleware)...)
}

func typeMiddleware(hooks []intake.MWHook) []func(http.Handler) http.Handler {
	var out []func(http.Handler) http.Handler

	for _, hook := range hooks {
		out = append(out, func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if hook.Timing == intake.PreHook {
					hook.Fn(eventFromRequest(r))
				}
				next.ServeHTTP(w, r)
				if hook.Timing == intake.PostHook {
					hook.Fn(eventFromRequest(r))
				}
			})
		})
	}

	return out
}

func eventFromRequest(r *http.Request) intake.Event {
	return intake.Event{
		Protocol: "http",
		Target:   r.URL.Path,
		ClientID: extractClientID(r),
	}
}

func extractClientID(r *http.Request) string {
	if v := r.Header.Get("X-Client-ID"); v != "" {
		return v
	}
	if v := r.Header.Get("X-API-Key"); v != "" {
		return v
	}
	return ""
}
