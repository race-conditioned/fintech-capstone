package policy

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/api_gateway/ports/outbound"
)

// Idempotency is a middleware that provides idempotency support for commands implementing IdempotentCommand.
func Idempotency[Com inbound.IdempotentCommand, Res inbound.Result](
	store outbound.Idempotency[Res],
	counterMetrics outbound.CounterMetrics,
) inbound.UnaryMiddleware[Com, Res] {
	return func(next inbound.UnaryHandler[Com, Res]) inbound.UnaryHandler[Com, Res] {
		return func(ctx context.Context, meta inbound.RequestMeta, cmd Com) (Res, error) {
			if store != nil {
				if cached, ok := store.Get(cmd.IdempotencyKey()); ok {
					if counterMetrics != nil {
						counterMetrics.IncIdempotentHit()
					}
					return cached, nil
				}
			}
			res, err := next(ctx, meta, cmd)
			if err == nil && store != nil && cmd.IdempotencyKey() != "" {
				store.Store(cmd.IdempotencyKey(), res)
			}
			return res, err
		}
	}
}
