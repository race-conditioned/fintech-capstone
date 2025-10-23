package grpc_transport

import (
	"context"
	pb "fintech-capstone/m/v2/internal/api_gateway/adapters/inbound/grpc/proto"
	"fintech-capstone/m/v2/internal/api_gateway/entrypoint"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
)

// TransferServer is the gRPC server for transfer operations.
type TransferServer struct {
	pb.UnimplementedTransferServiceServer
	h inbound.UnaryHandler[inbound.TransferCommand, inbound.TransferResult]
}

// NewTransferServer creates a new TransferServer.
func NewTransferServer(gw *entrypoint.Gateway) *TransferServer {
	return &TransferServer{h: gw.TransferHandler}
}

// Transfer handles transfer requests.
func (s *TransferServer) Transfer(ctx context.Context, req *pb.TransferCommand) (*pb.TransferResponse, error) {
	meta := metaFromGRPC(ctx, pb.TransferService_Transfer_FullMethodName)

	cmd := inbound.NewTransferCommand(
		req.GetFromAccount(),
		req.GetToAccount(),
		req.GetAmountCents(),
		req.GetIdempotencyKey(),
	)

	res, err := s.h(ctx, meta, cmd)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &pb.TransferResponse{
		TransactionId: res.TransactionID().String(),
		Status:        res.Status().String(),
		Message:       res.Message(),
	}, nil
}
