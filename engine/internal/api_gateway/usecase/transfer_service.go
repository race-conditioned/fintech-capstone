package usecase

import (
	"context"
	"errors"
	"fintech-capstone/m/v2/internal/api_gateway/apperr"
	"fintech-capstone/m/v2/internal/api_gateway/ports"
	"fintech-capstone/m/v2/internal/api_gateway/types"
)

type TransferService struct {
	dispatcher ports.Dispatcher
	metrics    ports.Metrics
}

func NewTransferService(d ports.Dispatcher, m ports.Metrics) *TransferService {
	return &TransferService{dispatcher: d, metrics: m}
}

func (s *TransferService) SubmitTransfer(ctx context.Context, meta ports.RequestMeta, cmd types.TransferCommand) (<-chan types.TransferResult, error) {
	if err := validate(cmd); err != nil {
		return nil, apperr.Invalid(err.Error())
	}
	// Delegate to worker pool via outbound port (transport-agnostic).
	return s.dispatcher.Submit(ctx, cmd), nil
}

func validate(cmd types.TransferCommand) error {
	if cmd.FromAccount == "" || cmd.ToAccount == "" {
		return errors.New("missing account IDs")
	}
	if cmd.AmountCents <= 0 {
		return errors.New("amount must be positive")
	}
	if cmd.IdempotencyKey == "" {
		return errors.New("missing idempotency key")
	}
	return nil
}
