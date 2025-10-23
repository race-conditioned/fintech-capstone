package outbound

import (
	"fintech-capstone/m/v2/internal/api_gateway/contracts"
	"time"
)

// Metrics defines counters and observations for telemetry.
type Metrics interface {
	IncRequest()
	IncSuccess()
	IncRateLimited()
	IncTimeout()
	IncIdempotentHit()
	ObserveLatency(d time.Duration)
	Snapshot() contracts.MetricsSnapshot
}
