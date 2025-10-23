package inbound

import (
	"context"

	"github.com/google/uuid"
)

// Inbound port: Application service entrypoint (called by HTTP/gRPC adapters).
type TransfersUseCase interface {
	SubmitTransfer(ctx context.Context, meta RequestMeta, cmd TransferCommand) (<-chan TransferResult, error)
}

// TransferCommand defines the external API payload for /transfer from any transport.
type TransferCommand struct {
	FromAccount    string
	ToAccount      string
	AmountCents    int64
	IdempotencyKey string
}

// Result returned by the use case.
// It is the internal representation of a completed transfer job.
type TransferResult struct {
	TransactionID uuid.UUID
	Status        string // "success" | "rejected" | "rate_limited" | "duplicate"
	Message       string
}
