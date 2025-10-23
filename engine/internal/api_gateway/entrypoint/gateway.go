package entrypoint

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/app/middleware"
	"fintech-capstone/m/v2/internal/api_gateway/app/policy"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/api_gateway/ports/outbound"
	"fintech-capstone/m/v2/internal/platform"
	"fmt"
	"time"
)

type Gateway struct {
	// composed, ready-to-mount handlers:
	transferH inbound.UnaryHandler[inbound.TransferCommand, inbound.TransferResult]

	// system dependencies for /metrics etc.
	metrics    outbound.Metrics
	dispatcher outbound.Dispatcher
	logger     platform.Logger
}

func NewGateway(
	transferUC inbound.TransfersUseCase,
	transferLimiter outbound.Limiter,
	transferIdemp outbound.Idempotency,
	systemMetrics outbound.Metrics,
	dispatcher outbound.Dispatcher,
	logger platform.Logger,
	transferTimeout time.Duration,
) *Gateway {
	base := func(ctx context.Context, meta inbound.RequestMeta, cmd inbound.TransferCommand) (inbound.TransferResult, error) {
		fmt.Println("base")
		ch, err := transferUC.SubmitTransfer(ctx, meta, cmd)
		if err != nil {
			return inbound.TransferResult{}, err
		}

		select {
		case <-ctx.Done():
			// Let the timeout middleware handle the metric increment on timeout
			return inbound.TransferResult{}, ctx.Err()
		case res := <-ch:
			return res, nil
		}
	}

	// Compose policy middlewares (order matters)
	// ---------------------
	// Idempotency policy is fired before rate limit policy as a design decision.
	// Retries are expected (client backoff + jitter).
	// When a duplicate key is used, the cached result is returned immediately, even if the client has hit other limits.
	// That reduces client churn and lowers total work across the system.
	// Idempotency lookup is fast & cheap (in-memory) and resilient.
	// Under a thundering herd of retries, this short-circuits work earlier than rate limit does.
	h := middleware.Chain[inbound.TransferCommand, inbound.TransferResult](
		base,
		middleware.CountRequests(systemMetrics),
		policy.IdempotencyTransfers(transferIdemp, systemMetrics),
		policy.RateLimitTransfers(transferLimiter, systemMetrics),
		policy.TimeoutTransfers(transferTimeout, systemMetrics),
		policy.ObserveLatencyTransfers(systemMetrics),
	)

	return &Gateway{
		transferH:  h,
		metrics:    systemMetrics,
		dispatcher: dispatcher,
		logger:     logger,
	}
}
