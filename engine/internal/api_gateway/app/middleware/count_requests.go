package middleware

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/api_gateway/ports/outbound"
	"fmt"
)

func CountRequests(m outbound.Metrics) inbound.UnaryMiddleware[inbound.TransferCommand, inbound.TransferResult] {
	return func(next inbound.UnaryHandler[inbound.TransferCommand, inbound.TransferResult]) inbound.UnaryHandler[inbound.TransferCommand, inbound.TransferResult] {
		return func(ctx context.Context, meta inbound.RequestMeta, cmd inbound.TransferCommand) (inbound.TransferResult, error) {
			fmt.Println("Count Requests")
			if m != nil {
				m.IncRequest()
			}
			return next(ctx, meta, cmd)
		}
	}
}
