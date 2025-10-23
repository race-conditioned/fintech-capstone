package policy

// PolicyStage represents different stages where policies can be applied.
type PolicyStage int

const (
	StageIdempotency PolicyStage = iota
	StageRateLimit
	StageTimeout
	StageLatency
)

// Policy represents a middleware policy at a specific stage.
type Policy struct {
	Stage PolicyStage
	MW    any
}

// DefaultPolicyOrder defines the order in which policies are applied.
// Idempotency policy is fired before rate limit policy as a design decision.
// Retries are expected (client backoff + jitter).
// When a duplicate key is used, the cached result is returned immediately, even if the client has hit other limits.
// That reduces client churn and lowers total work across the system.
// Idempotency lookup is fast & cheap (in-memory) and resilient.
// Under a thundering herd of retries, this short-circuits work earlier than rate limit does.
var DefaultPolicyOrder = []PolicyStage{
	StageIdempotency,
	StageRateLimit,
	StageTimeout,
	StageLatency,
}
