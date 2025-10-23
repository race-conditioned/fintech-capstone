package entrypoint

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
)

// TransferHandler handles transfer requests.
// It simply delegates to the configured middleware enriched handler in the gateway.
func (g *Gateway) TransferHandler(ctx context.Context, meta inbound.RequestMeta, cmd inbound.TransferCommand) (inbound.TransferResult, error) {
	return g.transferH(ctx, meta, cmd)
}
