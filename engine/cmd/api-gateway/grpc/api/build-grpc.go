package grpc_api

import (
	"fintech-capstone/m/v2/cmd/api-gateway/stubs"
	grpc_transport "fintech-capstone/m/v2/internal/api_gateway/adapters/inbound/grpc"
	pb "fintech-capstone/m/v2/internal/api_gateway/adapters/inbound/grpc/proto"
	"fintech-capstone/m/v2/internal/api_gateway/app"
	"fintech-capstone/m/v2/internal/api_gateway/entrypoint"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/platform"
	"fmt"
	"time"

	"google.golang.org/grpc"
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

	grpcSrv, err := grpc_transport.NewGRPCServer(":9090", func(gs *grpc.Server) {
		pb.RegisterTransferServiceServer(gs, grpc_transport.NewTransferServer(gw))
	})
	if err != nil {
		logger.Fatal(fmt.Errorf("grpc server init: %w", err))
	}
	return grpcSrv
}
