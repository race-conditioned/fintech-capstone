package gateway

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/ports"
	"fintech-capstone/m/v2/internal/api_gateway/types"
)

func (g *Gateway) TransferHandler(ctx context.Context, meta ports.RequestMeta, cmd types.TransferCommand) (types.TransferResult, error) {
	return g.transferH(ctx, meta, cmd)
}
