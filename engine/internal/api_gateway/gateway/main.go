package gateway

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/gateway/policy"
	"fintech-capstone/m/v2/internal/api_gateway/ports"
	"fintech-capstone/m/v2/internal/api_gateway/types"
	"time"
)

type Gateway struct {
	// composed, ready-to-mount handlers:
	transferH ports.UnaryHandler[types.TransferCommand, types.TransferResult]

	// system dependencies for /metrics etc.
	metrics    ports.Metrics
	dispatcher ports.Dispatcher
	logger     ports.Logger
}

func New(
	transferUC ports.TransfersUseCase,
	transferLimiter ports.Limiter,
	transferIdemp ports.Idempotency,
	systemMetrics ports.Metrics,
	dispatcher ports.Dispatcher,
	logger ports.Logger,
	transferTimeout time.Duration,
) *Gateway {
	base := func(ctx context.Context, meta ports.RequestMeta, cmd types.TransferCommand) (types.TransferResult, error) {
		if systemMetrics != nil {
			systemMetrics.IncRequest()
		}

		ch, err := transferUC.SubmitTransfer(ctx, meta, cmd)
		if err != nil {
			return types.TransferResult{}, err
		}

		select {
		case <-ctx.Done():
			// Let the timeout middleware handle the metric increment on timeout
			return types.TransferResult{}, ctx.Err()
		case res := <-ch:
			return res, nil
		}
	}

	// Compose policy middlewares (order matters)
	h := policy.Chain[types.TransferCommand, types.TransferResult](
		base,
		policy.ObserveLatencyTransfers(systemMetrics),
		policy.TimeoutTransfers(transferTimeout, systemMetrics),
		policy.IdempotencyTransfers(transferIdemp, systemMetrics),
		policy.RateLimitTransfers(transferLimiter, systemMetrics),
	)

	return &Gateway{
		transferH:  h,
		metrics:    systemMetrics,
		dispatcher: dispatcher,
		logger:     logger,
	}
}
