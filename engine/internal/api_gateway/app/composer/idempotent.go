package composer

import (
	"fintech-capstone/m/v2/internal/api_gateway/app/policy"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/api_gateway/ports/outbound"
)

// IdempotentComposer adds idempotency policy support.
type IdempotentComposer[Com inbound.IdempotentCommand, Res inbound.Result] struct {
	*Composer[Com, Res]
}

// Factory that returns a specialized composer for idempotent commands.
func NewIdempotentComposer[Com inbound.IdempotentCommand, Res inbound.Result](
	deps CommonDeps,
	idemp outbound.Idempotency[Res],
	metrics outbound.Metrics,
) *IdempotentComposer[Com, Res] {
	c := NewComposer[Com, Res](deps)
	c.idempMW = policy.Idempotency[Com, Res](idemp, metrics)
	return &IdempotentComposer[Com, Res]{Composer: c}
}
