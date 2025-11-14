package policy

import (
	"context"
	"time"

	"fintech-capstone/m/v2/internal/api_gateway/ports/outbound"

	hexa_inbound "github.com/race-conditioned/hexa/horizon/ports/inbound"
)

type AppHandler = hexa_inbound.UnaryHandler[Plugins, hexa_inbound.Command, hexa_inbound.Result]

type Plugins interface {
	context.Context
	Metrics() outbound.Metrics
	Limiter() outbound.Limiter
	Timeout() time.Duration
	Idempotency() outbound.Idempotency[hexa_inbound.Result]
}

type PluginsImpl struct {
	ctx         context.Context
	metrics     outbound.Metrics
	limiter     outbound.Limiter
	idempotency outbound.Idempotency[hexa_inbound.Result]
}

func NewPluginsImpl(
	ctx context.Context,
	metrics outbound.Metrics,
	limiter outbound.Limiter,
	idempotency outbound.Idempotency[hexa_inbound.Result],
) *PluginsImpl {
	return &PluginsImpl{
		ctx:         ctx,
		metrics:     metrics,
		limiter:     limiter,
		idempotency: idempotency,
	}
}

func (c *PluginsImpl) Metrics() outbound.Metrics { return c.metrics }
func (c *PluginsImpl) Limiter() outbound.Limiter { return c.limiter }
func (c *PluginsImpl) Timeout() time.Duration    { return 2 * time.Second }
func (c *PluginsImpl) Idempotency() outbound.Idempotency[hexa_inbound.Result] {
	return c.idempotency
}
func (p *PluginsImpl) Deadline() (time.Time, bool) { return p.ctx.Deadline() }
func (p *PluginsImpl) Done() <-chan struct{}       { return p.ctx.Done() }
func (p *PluginsImpl) Err() error                  { return p.ctx.Err() }
func (p *PluginsImpl) Value(key any) any           { return p.ctx.Value(key) }
