package grpc_api

import (
	"fintech-capstone/m/v2/cmd/api-gateway/stubs"
	grpc_transport "fintech-capstone/m/v2/internal/api_gateway/adapters/transports/grpc"
	pb "fintech-capstone/m/v2/internal/api_gateway/adapters/transports/grpc/proto"
	"fintech-capstone/m/v2/internal/api_gateway/gateway"
	"fintech-capstone/m/v2/internal/api_gateway/ports"
	"fintech-capstone/m/v2/internal/api_gateway/usecase"
	"fmt"
	"time"

	"google.golang.org/grpc"
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

	grpcSrv, err := grpc_transport.NewGRPCServer(":9090", func(gs *grpc.Server) {
		pb.RegisterTransferServiceServer(gs, grpc_transport.NewTransferServer(gw))
	})
	if err != nil {
		logger.Fatal(fmt.Errorf("grpc server init: %w", err))
	}
	return grpcSrv
}
