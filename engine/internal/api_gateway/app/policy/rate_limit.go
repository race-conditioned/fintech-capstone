package policy

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/api_gateway/ports/outbound"
	"fintech-capstone/m/v2/internal/platform/apperr"
)

// RateLimit is a middleware that enforces rate limiting based on the provided limiter.
func RateLimit[Com inbound.Command, Res inbound.Result](limiter outbound.Limiter, counterMetrics outbound.CounterMetrics) inbound.UnaryMiddleware[Com, Res] {
	return func(next inbound.UnaryHandler[Com, Res]) inbound.UnaryHandler[Com, Res] {
		return func(ctx context.Context, meta inbound.RequestMeta, cmd Com) (Res, error) {
			var zero Res
			if limiter != nil && !limiter.Allow(meta.ClientID) {
				if counterMetrics != nil {
					counterMetrics.IncRateLimited()
				}
				return zero, apperr.RateLimited("rate limit exceeded")
			}
			return next(ctx, meta, cmd)
		}
	}
}
