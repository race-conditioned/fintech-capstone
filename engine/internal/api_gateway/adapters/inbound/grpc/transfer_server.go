package grpc_transport

import (
	"context"
	pb "fintech-capstone/m/v2/internal/api_gateway/adapters/inbound/grpc/proto"
	"fintech-capstone/m/v2/internal/api_gateway/entrypoint"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
)

type TransferServer struct {
	pb.UnimplementedTransferServiceServer
	h inbound.UnaryHandler[inbound.TransferCommand, inbound.TransferResult]
}

func NewTransferServer(gw *entrypoint.Gateway) *TransferServer {
	return &TransferServer{h: gw.TransferHandler}
}

func (s *TransferServer) Transfer(ctx context.Context, req *pb.TransferCommand) (*pb.TransferResponse, error) {
	meta := metaFromGRPC(ctx, pb.TransferService_Transfer_FullMethodName)

	cmd := inbound.TransferCommand{
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
