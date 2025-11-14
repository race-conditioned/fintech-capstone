package policy

import (
	"fmt"

	"fintech-capstone/m/v2/internal/platform/apperr"

	hexa_inbound "github.com/race-conditioned/hexa/horizon/ports/inbound"
)

// // RateLimit is a middleware that enforces rate limiting based on the provided limiter.
// func RateLimit(next inbound.UnaryHandler[Com, Res]) inbound.UnaryHandler[Com, Res] {
// 	return func(ctx context.Context, meta inbound.RequestMeta, cmd Com) (Res, error) {
// 		var zero Res
// 		if limiter != nil && !limiter.Allow(meta.ClientID) {
// 			if counterMetrics != nil {
// 				counterMetrics.IncRateLimited()
// 			}
// 			return zero, apperr.RateLimited("rate limit exceeded")
// 		}
// 		return next(ctx, meta, cmd)
// 	}
// }

// RateLimit is a middleware that enforces rate limiting based on the provided limiter.
func RateLimit(next AppHandler) AppHandler {
	return func(ctx Plugins, meta hexa_inbound.RequestMeta, cmd hexa_inbound.Command) (hexa_inbound.Result, error) {
		fmt.Println("Applying rate limit...")
		var zero hexa_inbound.Result
		if !ctx.Limiter().Allow(meta.ClientID) {
			ctx.Metrics().IncRateLimited()
			return zero, apperr.RateLimited("rate limit exceeded")
		}
		return next(ctx, meta, cmd)
	}
}
