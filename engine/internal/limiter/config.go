package limiter

import "time"

// PerClientConfig controls per-client token buckets.
type PerClientConfig struct {
	// RatePerSec is tokens added per second. 0 or negative => unlimited.
	RatePerSec float64
	// Burst is the maximum bucket size (must be >= 1 when RatePerSec > 0).
	Burst int
	// InitialTokens seeds new client buckets. If <= 0, defaults to Burst.
	InitialTokens float64
	// TTL evicts inactive clients after this duration. If 0, defaults to 10m.
	TTL time.Duration
}

// GlobalConfig controls the optional global token bucket.
type GlobalConfig struct {
	RatePerSec float64
	Burst      int
	// InitialTokens seeds the bucket. If <= 0, defaults to Burst.
	InitialTokens float64
}

// Config combines all options.
type Config struct {
	PerClient       PerClientConfig
	Global          *GlobalConfig // nil => no global limiting
	NumShards       int           // 0 => default (64)
	CleanupInterval time.Duration // 0 => default (1m)
	Clock           Clock         // nil => system clock
}
