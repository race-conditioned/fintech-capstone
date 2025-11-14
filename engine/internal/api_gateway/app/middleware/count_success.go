package middleware

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/api_gateway/ports/outbound"
)

// CountSuccess is a middleware that counts incoming requests using the provided metrics.
func CountSuccess[Com inbound.Command, Result inbound.Result](counterMetrics outbound.CounterMetrics) inbound.UnaryMiddleware[Com, Result] {
	return func(next inbound.UnaryHandler[Com, Result]) inbound.UnaryHandler[Com, Result] {
		return func(ctx context.Context, meta inbound.RequestMeta, cmd Com) (Result, error) {
			r, err := next(ctx, meta, cmd)
			if err == nil && counterMetrics != nil {
				counterMetrics.IncSuccess()
			}
			return r, err
		}
	}
}
