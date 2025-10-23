// internal/api_gateway/app/composer/composer.go
package composer

import (
	"fintech-capstone/m/v2/internal/api_gateway/app/middleware"
	"fintech-capstone/m/v2/internal/api_gateway/app/policy"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/api_gateway/ports/outbound"
	"time"
)

// shared dependency preset
type CommonDeps struct {
	Metrics outbound.Metrics
	Limiter outbound.Limiter
	Timeout time.Duration
}

// Composer wraps middleware around a handler
type Composer[Com inbound.Command, Res inbound.Result] struct {
	metrics   outbound.Metrics
	limiter   outbound.Limiter
	timeout   time.Duration
	idempMW   inbound.UnaryMiddleware[Com, Res]
	pre, post []inbound.UnaryMiddleware[Com, Res]
}

// Option configures a Composer
type Option[Com inbound.Command, Res inbound.Result] func(*Composer[Com, Res])

// WithMetrics adds metrics middleware
func (c *Composer[Com, Res]) WithMetrics(m outbound.Metrics) Option[Com, Res] {
	return func(c *Composer[Com, Res]) { c.metrics = m }
}

// WithLimiter adds limiter middleware
func (c *Composer[Com, Res]) WithLimiter(l outbound.Limiter) Option[Com, Res] {
	return func(c *Composer[Com, Res]) { c.limiter = l }
}

// WithTimeout adds timeout middleware
func (c *Composer[Com, Res]) WithTimeout(d time.Duration) Option[Com, Res] {
	return func(c *Composer[Com, Res]) { c.timeout = d }
}

// NewComposer creates a Composer with common middlewares configured from deps
func NewComposer[Com inbound.Command, Res inbound.Result](deps CommonDeps) *Composer[Com, Res] {
	c := &Composer[Com, Res]{}
	for _, opt := range []Option[Com, Res]{
		c.WithMetrics(deps.Metrics),
		c.WithLimiter(deps.Limiter),
		c.WithTimeout(deps.Timeout),
	} {
		opt(c)
	}
	return c
}

// Build constructs the final handler by composing middlewares around the base handler
func (c *Composer[Com, Res]) Build(
	base inbound.UnaryHandler[Com, Res],
	extra ...inbound.UnaryMiddleware[Com, Res],
) inbound.UnaryHandler[Com, Res] {
	policies := map[policy.PolicyStage]inbound.UnaryMiddleware[Com, Res]{}

	if c.idempMW != nil {
		policies[policy.StageIdempotency] = c.idempMW
	}
	if c.limiter != nil {
		policies[policy.StageRateLimit] = policy.RateLimit[Com, Res](c.limiter, c.metrics)
	}
	if c.timeout > 0 {
		policies[policy.StageTimeout] = policy.Timeout[Com, Res](c.timeout, c.metrics)
	}
	if c.metrics != nil {
		policies[policy.StageLatency] = policy.ObserveLatency[Com, Res](c.metrics)
	}

	var chain []inbound.UnaryMiddleware[Com, Res]

	chain = append(chain, middleware.CountRequests[Com, Res](c.metrics))

	for _, stage := range policy.DefaultPolicyOrder {
		if mw, ok := policies[stage]; ok {
			chain = append(chain, mw)
		}
	}

	for _, extraMW := range extra {
		chain = append(chain, extraMW)
	}

	chain = append(chain, middleware.CountSuccess[Com, Res](c.metrics))

	return middleware.Chain(base, chain...) // Compose policy middlewares (order matters)
}
