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
	fromAccount    string
	toAccount      string
	amountCents    int64
	idempotencyKey string
}

// NewTransferCommand creates a new TransferCommand.
func NewTransferCommand(fromAccount, toAccount string, amountCents int64, idempotencyKey string) TransferCommand {
	return TransferCommand{
		fromAccount:    fromAccount,
		toAccount:      toAccount,
		amountCents:    amountCents,
		idempotencyKey: idempotencyKey,
	}
}

// FromAccount returns the source account ID.
func (t TransferCommand) FromAccount() string {
	return t.fromAccount
}

// ToAccount returns the destination account ID.
func (t TransferCommand) ToAccount() string {
	return t.toAccount
}

// AmountCents returns the transfer amount in cents.
func (t TransferCommand) AmountCents() int64 {
	return t.amountCents
}

// IdempotencyKey returns the idempotency key for the transfer.
func (t TransferCommand) IdempotencyKey() string {
	return t.idempotencyKey
}

// TransferResult returned by the use case.
// It is the internal representation of a completed transfer job.
type TransferResult struct {
	transactionID uuid.UUID
	status        ResultStatus
	message       string
}

// NewTransferResult creates a new TransferResult.
func NewTransferResult(transactionID uuid.UUID, status ResultStatus, message string) TransferResult {
	return TransferResult{
		transactionID: transactionID,
		status:        status,
		message:       message,
	}
}

// TransactionID returns the transaction ID of the transfer.
func (t TransferResult) TransactionID() uuid.UUID {
	return t.transactionID
}

// Status returns the status of the transfer.
func (t TransferResult) Status() ResultStatus {
	return t.status
}

// Message returns the message associated with the transfer result.
func (t TransferResult) Message() string {
	return t.message
}
