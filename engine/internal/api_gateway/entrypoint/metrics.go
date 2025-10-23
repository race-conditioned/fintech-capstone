package entrypoint

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/contracts"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
)

func (g *Gateway) MetricsHandler(ctx context.Context, _ inbound.RequestMeta, _ struct{}) (contracts.MetricsSnapshot, error) {
	s := g.metrics.Snapshot()
	s.ActiveWorkers = g.dispatcher.ActiveWorkers()
	s.QueueDepth = g.dispatcher.QueueDepth()
	return s, nil
}
