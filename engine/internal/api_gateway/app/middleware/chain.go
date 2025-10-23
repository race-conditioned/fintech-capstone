package middleware

import "fintech-capstone/m/v2/internal/api_gateway/ports/inbound"

// Chain applies middlewares inside-out (last wraps closest to the handler).
func Chain[Com inbound.Command, Result inbound.Result](h inbound.UnaryHandler[Com, Result], mws ...inbound.UnaryMiddleware[Com, Result]) inbound.UnaryHandler[Com, Result] {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
