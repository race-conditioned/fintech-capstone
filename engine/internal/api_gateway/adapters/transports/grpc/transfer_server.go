package grpc_transport

import (
	"context"
	pb "fintech-capstone/m/v2/internal/api_gateway/adapters/transports/grpc/proto"
	"fintech-capstone/m/v2/internal/api_gateway/gateway"
	"fintech-capstone/m/v2/internal/api_gateway/ports"
	"fintech-capstone/m/v2/internal/api_gateway/types"
)

type TransferServer struct {
	pb.UnimplementedTransferServiceServer
	h ports.UnaryHandler[types.TransferCommand, types.TransferResult]
}

func NewTransferServer(gw *gateway.Gateway) *TransferServer {
	return &TransferServer{h: gw.TransferHandler}
}

func (s *TransferServer) Transfer(ctx context.Context, req *pb.TransferCommand) (*pb.TransferResponse, error) {
	meta := metaFromGRPC(ctx, pb.TransferService_Transfer_FullMethodName)

	cmd := types.TransferCommand{
		FromAccount:    req.GetFromAccount(),
		ToAccount:      req.GetToAccount(),
		AmountCents:    req.GetAmountCents(),
		IdempotencyKey: req.GetIdempotencyKey(),
	}

	res, err := s.h(ctx, meta, cmd)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &pb.TransferResponse{
		TransactionId: res.TransactionID.String(),
		Status:        res.Status,
		Message:       res.Message,
	}, nil
}
