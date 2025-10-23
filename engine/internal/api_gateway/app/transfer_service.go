package app

import (
	"context"
	"errors"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/api_gateway/ports/outbound"
	"fintech-capstone/m/v2/internal/platform"
	"fintech-capstone/m/v2/internal/platform/apperr"
)

type TransferService struct {
	dispatcher outbound.Dispatcher
	metrics    outbound.Metrics
	logger     platform.Logger
}

func NewTransferService(d outbound.Dispatcher, m outbound.Metrics, l platform.Logger) *TransferService {
	return &TransferService{dispatcher: d, metrics: m, logger: l}
}

func (s *TransferService) SubmitTransfer(ctx context.Context, meta inbound.RequestMeta, cmd inbound.TransferCommand) (<-chan inbound.TransferResult, error) {
	if err := validate(cmd); err != nil {
		return nil, apperr.Invalid(err.Error())
	}
	// Delegate to worker pool via outbound port (transport-agnostic).
	return s.dispatcher.Submit(ctx, cmd), nil
}

func validate(cmd inbound.TransferCommand) error {
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
