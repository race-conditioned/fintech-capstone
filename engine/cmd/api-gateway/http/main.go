package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"fintech-capstone/m/v2/cmd/api-gateway/stubs"
	"fintech-capstone/m/v2/internal/api_gateway/app"
	"fintech-capstone/m/v2/internal/api_gateway/app/policy"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/limiter"
	zap_adapter "fintech-capstone/m/v2/internal/platform/adapters/zap"

	"github.com/race-conditioned/hexa/endurance"
	"github.com/race-conditioned/hexa/fusion/dt"
	"github.com/race-conditioned/hexa/fusion/intake"
	"github.com/race-conditioned/hexa/horizon"
	"github.com/race-conditioned/hexa/symphony"
	"go.uber.org/zap"
)

// main is the entry point for the API Gateway HTTP server.
func main() {
	zaplog, err := zap.NewProduction()
	if err != nil {
		log.Println(fmt.Errorf("Init zap logger: %w", err).Error())
		return
	}
	defer zaplog.Sync()

	logger := zap_adapter.New(zaplog)

	// httpSrv := http_api.BuildServer(logger)

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

	uc := app.NewTransferService(dispatch, metrics, logger)
	plugins := policy.NewPluginsImpl(
		context.Background(),
		metrics,
		lim,
		idemp,
	)
	gw := horizon.NewGateway[policy.Plugins](plugins)

	// 	var DefaultPolicyOrder = []PolicyStage{
	// 	StageIdempotency,
	// 	StageRateLimit,
	// 	StageTimeout,
	// 	StageLatency,
	// }
	//

	idempotency := symphony.PolicyStage("idempotency")
	rateLimit := symphony.PolicyStage("rate_limit")
	timeout := symphony.PolicyStage("timeout")
	latency := symphony.PolicyStage("latency")

	composer := symphony.New[policy.Plugins](symphony.Order(), symphony.Order())

	DefaultPolicyOrder := []symphony.PolicyStage{
		"idempotency",
		"rate_limit",
		"timeout",
		"latency",
	}
	mid := symphony.Order(DefaultPolicyOrder...)

	transferComposition := symphony.Compose(
		composer,
		mid,
		symphony.WithPolicy(rateLimit, symphony.Lift[policy.Plugins, inbound.TransferCommand, inbound.TransferResult](policy.RateLimit)),
		symphony.WithPolicy(timeout, symphony.Lift[policy.Plugins, inbound.TransferCommand, inbound.TransferResult](policy.Timeout)),
		symphony.WithPolicy(latency, symphony.Lift[policy.Plugins, inbound.TransferCommand, inbound.TransferResult](policy.ObserveLatency)),
		symphony.WithPolicy(idempotency, symphony.LiftCap[policy.Plugins, inbound.TransferCommand, inbound.TransferResult](policy.Idempotency)),
	)

	h := transferComposition.Wrap(endurance.Transport(uc.SubmitTransfer, nil, nil))

	gw.RegisterHandler("transfer", horizon.Adapt(h))

	spec := intake.Spec{}

	routes := []dt.Route[policy.Plugins]{
		jsonRoute[inbound.TransferCommandHTTP]("transfer"),
	}

	fusion := dt.NewFusion[policy.Plugins](plugins, spec, gw, routes)
	router := fusion.Build()

	httpSrv, err := dt.NewHTTPServer(
		":8080", router,
		2*time.Second,  // ReadHeaderTimeout
		5*time.Second,  // ReadTimeout
		10*time.Second, // WriteTimeout
		60*time.Second, // IdleTimeout
	)
	if err != nil {
		log.Fatal(fmt.Errorf("http server init: %w", err))
	}

	// Context + signals for graceful shutdown
	httpCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start
	go func() {
		if err := httpSrv.Start(httpCtx); err != nil {
			log.Fatal(fmt.Errorf("running server: %w", err)) // platform.Field{Key: "component", Value: "http"})
			stop()
		}
	}()

	<-httpCtx.Done()

	// Shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(shutdownCtx)

	// // [policy.Plugins, inbound.IdempotentCommand, inbound.TransferResult]
	//
	// // Context + signals for graceful shutdown
	// ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// defer stop()
	//
	// // Start
	// go func() {
	// 	if err := httpSrv.Start(ctx); err != nil {
	// 		logger.Error(fmt.Errorf("running server: %w", err), platform.Field{Key: "component", Value: "http"})
	// 		stop()
	// 	}
	// }()
	//
	// <-ctx.Done()
	//
	// Shutdown
	// shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	// defer cancel()
	// _ = httpSrv.Shutdown(shutdownCtx)
}

func jsonRoute[T any](key string) dt.Route[policy.Plugins] {
	return dt.JSONRoute[policy.Plugins, T](horizon.HandlerKey(key))
}
