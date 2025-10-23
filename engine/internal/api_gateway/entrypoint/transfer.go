package entrypoint

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fmt"
)

func (g *Gateway) TransferHandler(ctx context.Context, meta inbound.RequestMeta, cmd inbound.TransferCommand) (inbound.TransferResult, error) {
	fmt.Println("transfer handler")
	return g.transferH(ctx, meta, cmd)
}
