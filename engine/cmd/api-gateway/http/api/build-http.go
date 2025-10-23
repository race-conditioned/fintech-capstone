package http_api

import (
	"fintech-capstone/m/v2/cmd/api-gateway/gateway"
	http_transport "fintech-capstone/m/v2/internal/api_gateway/adapters/inbound/http"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/platform"
	"fmt"
	"time"
)

// BuildServer builds and returns an HTTP server for the API Gateway.
func BuildServer(logger platform.Logger) inbound.Server {
	handler := http_transport.NewRouter(gateway.BuildGateway(logger), logger)

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
