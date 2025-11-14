package gateway

import (
	"context"
	"time"

	"fintech-capstone/m/v2/cmd/api-gateway/stubs"
	"fintech-capstone/m/v2/internal/api_gateway/app"
	"fintech-capstone/m/v2/internal/api_gateway/app/composer"
	"fintech-capstone/m/v2/internal/api_gateway/entrypoint"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/limiter"
	"fintech-capstone/m/v2/internal/platform"
)

// BuildGateway constructs the API Gateway with all its handlers and dependencies.
// It is customised for our application domain (transfers) with gateway level middleware:
// Metrics, Limiting, Idempotency, Timeout.
func BuildGateway(logger platform.Logger) *entrypoint.Gateway {
	_, idemp, dispatch, metrics := stubs.BuildTransfer()

	lim := limiter.New(context.Background(), limiter.Config{
		PerClient: limiter.PerClientConfig{
			RatePerSec:    50,  // 50 rps per client
			Burst:         100, // allow short spikes
			InitialTokens: 100, // warm start
			TTL:           15 * time.Minute,
		},
		Global: &limiter.GlobalConfig{
			RatePerSec:    1000, // cap overall ingress
			Burst:         2000,
			InitialTokens: 2000,
		},
		NumShards:       64,
		CleanupInterval: time.Minute,
	})
	// Use case (app layer)
	uc := app.NewTransferService(dispatch, metrics, logger)

	// Endpoints provider (base handlers only)

	deps := composer.CommonDeps{
		Metrics: metrics,
		Limiter: lim,
		Timeout: 2 * time.Second,
	}

	compTR := composer.NewIdempotentComposer[inbound.TransferCommand, inbound.TransferResult](deps, idemp)

	// Build composed handlers per endpoint, no repeated options
	submitH := compTR.Build(uc.SubmitTransfer)

	// Mount on gateway (kept dumb)
	gw := entrypoint.NewGateway(metrics, dispatch, logger,
		entrypoint.WithTransfer(submitH),
		// entrypoint.WithTransferCancel(cancelH), - example more endpoints
	)
	return gw
}
