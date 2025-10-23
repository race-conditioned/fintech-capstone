package policy

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/api_gateway/ports/outbound"
	"time"
)

// ObserveLatency is a middleware that observes the latency of requests and records success metrics.
func ObserveLatency[Com inbound.Command, Result inbound.Result](latencyMetrics outbound.LatencyMetrics) inbound.UnaryMiddleware[Com, Result] {
	return func(next inbound.UnaryHandler[Com, Result]) inbound.UnaryHandler[Com, Result] {
		return func(ctx context.Context, meta inbound.RequestMeta, cmd Com) (Result, error) {
			start := time.Now()
			res, err := next(ctx, meta, cmd)
			if latencyMetrics != nil {
				latencyMetrics.ObserveLatency(time.Since(start))
			}
			return res, err
		}
	}
}
