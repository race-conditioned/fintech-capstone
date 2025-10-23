package outbound

import (
	"fintech-capstone/m/v2/internal/api_gateway/contracts"
	"time"
)

// Metrics aggregates all metric capabilities.
// Existing code can still depend on this for convenience.
type Metrics interface {
	CounterMetrics
	LatencyMetrics
	SnapshotMetrics
}

// CounterMetrics defines event counters.
type CounterMetrics interface {
	IncRequest()
	IncSuccess()
	IncRateLimited()
	IncTimeout()
	IncIdempotentHit()
}

// LatencyMetrics defines latency observation.
type LatencyMetrics interface {
	ObserveLatency(d time.Duration)
}

// SnapshotMetrics defines exportable snapshotting.
type SnapshotMetrics interface {
	Snapshot() contracts.MetricsSnapshot
}
