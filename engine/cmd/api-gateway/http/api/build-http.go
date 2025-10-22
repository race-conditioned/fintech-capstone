package http_api

import (
	"fintech-capstone/m/v2/cmd/api-gateway/stubs"
	http_transport "fintech-capstone/m/v2/internal/api_gateway/adapters/transports/http"
	"fintech-capstone/m/v2/internal/api_gateway/gateway"
	"fintech-capstone/m/v2/internal/api_gateway/ports"
	"fintech-capstone/m/v2/internal/api_gateway/usecase"
	"fmt"
	"time"
)

func BuildServer(logger ports.Logger) ports.InboundServer {
	limiter, idemp, dispatch, metrics := stubs.Build()

	uc := usecase.NewTransferService(dispatch, metrics)

	transferLimiter := limiter
	transferIdemp := idemp

	gw := gateway.New(
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
		5*time.Second,  // ReadHeaderTimeout
		10*time.Second, // ReadTimeout
		30*time.Second, // WriteTimeout
		60*time.Second, // IdleTimeout
	)
	if err != nil {
		logger.Fatal(fmt.Errorf("http server init: %w", err))
	}
	return httpSrv
}
