package entrypoint

import (
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/api_gateway/ports/outbound"
	"fintech-capstone/m/v2/internal/platform"
)

// Gateway is the API Gateway entrypoint, composing handlers with middleware.
type Gateway struct {
	transferH  inbound.UnaryHandler[inbound.TransferCommand, inbound.TransferResult]
	metrics    outbound.Metrics
	dispatcher outbound.Dispatcher
	logger     platform.Logger
}

// Option configures a Gateway.
type Option func(*Gateway)

// WithTransfer sets the transfer handler.
func WithTransfer(h inbound.UnaryHandler[inbound.TransferCommand, inbound.TransferResult]) Option {
	return func(g *Gateway) { g.transferH = h }
}

// NewGateway constructs a new API Gateway entrypoint with all handlers composed with middleware.
func NewGateway(
	metrics outbound.Metrics,
	dispatcher outbound.Dispatcher,
	logger platform.Logger,
	opts ...Option,
) *Gateway {
	g := &Gateway{
		metrics:    metrics,
		dispatcher: dispatcher,
		logger:     logger,
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}
