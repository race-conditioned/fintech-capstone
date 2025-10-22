package ports

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/types"
)

// Inbound port: Application service entrypoint (called by HTTP/gRPC adapters).
type TransfersUseCase interface {
	SubmitTransfer(ctx context.Context, meta RequestMeta, cmd types.TransferCommand) (<-chan types.TransferResult, error)
}

type RequestMeta struct {
	ClientID  string
	RequestID string
	TraceID   string
	RemoteIP  string
	Protocol  string // "http","grpc",...
	Target    string // path or method name for logging
}

type (
	UnaryHandler[Req any, Res any]    func(ctx context.Context, meta RequestMeta, req Req) (Res, error)
	UnaryMiddleware[Req any, Res any] func(next UnaryHandler[Req, Res]) UnaryHandler[Req, Res]
)
