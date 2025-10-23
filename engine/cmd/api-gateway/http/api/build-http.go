package http_api

import (
	"fintech-capstone/m/v2/cmd/api-gateway/stubs"
	http_transport "fintech-capstone/m/v2/internal/api_gateway/adapters/inbound/http"
	"fintech-capstone/m/v2/internal/api_gateway/app"
	"fintech-capstone/m/v2/internal/api_gateway/entrypoint"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/platform"
	"fmt"
	"time"
)

func BuildServer(logger platform.Logger) inbound.Server {
	limiter, idemp, dispatch, metrics := stubs.Build()

	uc := app.NewTransferService(dispatch, metrics, logger)

	transferLimiter := limiter
	transferIdemp := idemp

	gw := entrypoint.NewGateway(
		uc,
		transferLimiter,
		transferIdemp,
		metrics,
		dispatch,
		logger,
		2*time.Second,
	)

	// Build HTTP handler and server
	handler := http_transport.NewRouter(gw, logger)

	httpSrv, err := http_transport.NewHTTPServer(
		":8080", handler,
		2*time.Second,  // ReadHeaderTimeout
		5*time.Second,  // ReadTimeout
		10*time.Second, // WriteTimeout
		60*time.Second, // IdleTimeout
	)
	if err != nil {
		logger.Fatal(fmt.Errorf("http server init: %w", err))
	}
	return httpSrv
}
