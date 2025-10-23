package grpc_api

import (
	"fintech-capstone/m/v2/cmd/api-gateway/gateway"
	grpc_transport "fintech-capstone/m/v2/internal/api_gateway/adapters/inbound/grpc"
	pb "fintech-capstone/m/v2/internal/api_gateway/adapters/inbound/grpc/proto"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/platform"
	"fmt"

	"google.golang.org/grpc"
)

// BuildServer builds and returns a gRPC server for the API Gateway.
func BuildServer(logger platform.Logger) inbound.Server {
	grpcSrv, err := grpc_transport.NewGRPCServer(":9090", func(gs *grpc.Server) {
		pb.RegisterTransferServiceServer(gs, grpc_transport.NewTransferServer(gateway.BuildGateway(logger)))
	})
	if err != nil {
		logger.Fatal(fmt.Errorf("grpc server init: %w", err))
	}
	return grpcSrv
}
