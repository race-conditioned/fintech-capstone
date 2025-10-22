// cmd/api-gateway/main.go
package main

import (
	"context"
	"errors"
	grpc_api "fintech-capstone/m/v2/cmd/api-gateway/grpc/api"
	http_api "fintech-capstone/m/v2/cmd/api-gateway/http/api"
	zap_adapter "fintech-capstone/m/v2/internal/api_gateway/adapters/loggers/zap"
	"fintech-capstone/m/v2/internal/api_gateway/ports"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func runAll(ctx context.Context, servers ...ports.InboundServer) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, s := range servers {
		srv := s
		g.Go(func() error { return srv.Start(ctx) })
	}

	// Wait for shutdown signal
	<-ctx.Done()

	// Create shutdown context
	stop, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var shutdownErr error
	for _, s := range servers {
		if err := s.Shutdown(stop); err != nil {
			shutdownErr = multierr.Append(shutdownErr, fmt.Errorf("%T shutdown: %w", s, err))
		}
	}

	// Wait for all goroutines to finish
	runErr := g.Wait()

	// If shutdown was triggered via context cancel, that's not fatal
	if errors.Is(ctx.Err(), context.Canceled) {
		ctx = context.Background() // ignore canceled context
		runErr = nil
	}

	// Combine runErr + shutdownErr
	return multierr.Combine(runErr, shutdownErr)
}

func main() {
	zaplog, err := zap.NewProduction()
	if err != nil {
		log.Println(fmt.Errorf("Init zap logger: %w", err).Error())
		return
	}
	defer func() { _ = zaplog.Sync() }()

	logger := zap_adapter.New(zaplog)
	httpSrv := http_api.BuildServer(logger)
	grpcSrv := grpc_api.BuildServer(logger)

	// build gw, httpSrv, grpcSrv...
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	_ = runAll(ctx, httpSrv, grpcSrv)
}
