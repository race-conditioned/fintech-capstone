package endpoints

import (
	"context"
	"errors"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/platform/apperr"
)

// ProviderTransfers wraps use case to provide endpoint handlers.
type ProviderTransfers struct {
	uc inbound.TransfersUseCase
}

// NewProviderTransfers creates a new ProviderTransfers.
func NewProviderTransfers(uc inbound.TransfersUseCase) *ProviderTransfers {
	return &ProviderTransfers{uc: uc}
}

// SubmitBase handles transfer submission.
func (p *ProviderTransfers) SubmitBase() inbound.UnaryHandler[inbound.TransferCommand, inbound.TransferResult] {
	return func(ctx context.Context, meta inbound.RequestMeta, cmd inbound.TransferCommand) (inbound.TransferResult, error) {
		ch, err := p.uc.SubmitTransfer(ctx, meta, cmd)
		if err != nil {
			return inbound.TransferResult{}, err
		}

		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return inbound.TransferResult{}, apperr.Timeout("processing timeout")
			}
			return inbound.TransferResult{}, apperr.Internal("request canceled")

		case res := <-ch:
			return res, nil
		}
	}
}
