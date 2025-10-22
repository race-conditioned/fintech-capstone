package gateway

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/ports"
	"fintech-capstone/m/v2/internal/api_gateway/types"
)

func (g *Gateway) MetricsHandler(ctx context.Context, _ ports.RequestMeta, _ struct{}) (types.MetricsSnapshot, error) {
	s := g.metrics.Snapshot()
	s.ActiveWorkers = g.dispatcher.ActiveWorkers()
	s.QueueDepth = g.dispatcher.QueueDepth()
	return s, nil
}
