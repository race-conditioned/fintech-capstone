// cmd/api-gateway-http/main.go
package main

import (
	"context"
	http_api "fintech-capstone/m/v2/cmd/api-gateway/http/api"
	zap_adapter "fintech-capstone/m/v2/internal/api_gateway/adapters/loggers/zap"
	"fintech-capstone/m/v2/internal/api_gateway/ports"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	zaplog, err := zap.NewProduction()
	if err != nil {
		log.Println(fmt.Errorf("Init zap logger: %w", err).Error())
		return
	}
	defer zaplog.Sync()

	logger := zap_adapter.New(zaplog)

	httpSrv := http_api.BuildServer(logger)

	// Context + signals for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start
	go func() {
		if err := httpSrv.Start(ctx); err != nil {
			logger.Error(fmt.Errorf("running server: %w", err), ports.Field{Key: "component", Value: "http"})
			stop()
		}
	}()

	<-ctx.Done()

	// Shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(shutdownCtx)
}
