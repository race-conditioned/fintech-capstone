package gateway

import (
	"fintech-capstone/m/v2/cmd/api-gateway/stubs"
	"fintech-capstone/m/v2/internal/api_gateway/app"
	"fintech-capstone/m/v2/internal/api_gateway/app/composer"
	"fintech-capstone/m/v2/internal/api_gateway/app/endpoints"
	"fintech-capstone/m/v2/internal/api_gateway/entrypoint"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/platform"
	"time"
)

// BuildGateway constructs the API Gateway with all its handlers and dependencies.
// It is customised for our application domain (transfers) with gateway level middleware:
// Metrics, Limiting, Idempotency, Timeout.
func BuildGateway(logger platform.Logger) *entrypoint.Gateway {
	limiter, idemp, dispatch, metrics := stubs.BuildTransfer()

	// Use case (app layer)
	uc := app.NewTransferService(dispatch, metrics, logger)

	// Endpoints provider (base handlers only)
	trEP := endpoints.NewProviderTransfers(uc)

	deps := composer.CommonDeps{
		Metrics: metrics,
		Limiter: limiter,
		Timeout: 2 * time.Second,
	}

	compTR := composer.NewIdempotentComposer[inbound.TransferCommand, inbound.TransferResult](deps, idemp, metrics)

	// Build composed handlers per endpoint, no repeated options
	submitH := compTR.Build(trEP.SubmitBase())

	// Mount on gateway (kept dumb)
	gw := entrypoint.NewGateway(metrics, dispatch, logger,
		entrypoint.WithTransfer(submitH),
		// entrypoint.WithTransferCancel(cancelH),
	)
	return gw
}
