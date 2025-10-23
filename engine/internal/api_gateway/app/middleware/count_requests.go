package middleware

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/api_gateway/ports/outbound"
)

// CountRequests is a middleware that counts incoming requests using the provided metrics.
func CountRequests[Com inbound.Command, Result inbound.Result](m outbound.Metrics) inbound.UnaryMiddleware[Com, Result] {
	return func(next inbound.UnaryHandler[Com, Result]) inbound.UnaryHandler[Com, Result] {
		return func(ctx context.Context, meta inbound.RequestMeta, cmd Com) (Result, error) {
			if m != nil {
				m.IncRequest()
			}
			return next(ctx, meta, cmd)
		}
	}
}
