package symphony

// PolicyStage labels a named slot in the pipeline (idempotency, rate_limit, timeout, latency, etc.)
type PolicyStage string

// PolicyOrder defines a deterministic order for stages.
// Using a type with Order() enforces configuration at compile-time (no empty defaults).
type PolicyOrder interface {
	Order() []PolicyStage
}

type staticOrder []PolicyStage

func (o staticOrder) Order() []PolicyStage { return o }

func Order(stages ...PolicyStage) PolicyOrder {
	return staticOrder(stages)
}
