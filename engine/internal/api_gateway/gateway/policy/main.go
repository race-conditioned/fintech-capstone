package policy

import (
	"fintech-capstone/m/v2/internal/api_gateway/ports"
)

// Chain applies middlewares inside-out (last wraps closest to the handler).
func Chain[Req any, Res any](h ports.UnaryHandler[Req, Res], mws ...ports.UnaryMiddleware[Req, Res]) ports.UnaryHandler[Req, Res] {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
