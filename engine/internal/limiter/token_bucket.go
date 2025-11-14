package limiter

import (
	"math"
	"sync"
	"time"
)

// tokenBucket is a classic token-bucket with lazy refill.
// All methods require mu to be held by the caller.
type tokenBucket struct {
	mu         sync.Mutex
	tokens     float64   // current token count
	lastRefill time.Time // last refill timestamp (monotonic)
	rate       float64   // tokens/sec
	burst      float64   // max tokens
}

func newTokenBucket(now time.Time, rate float64, burst int, initial float64) *tokenBucket {
	b := &tokenBucket{
		tokens:     clampInit(initial, burst),
		lastRefill: now,
		rate:       rate,
		burst:      float64(max(1, burst)),
	}
	if rate <= 0 {
		// Unlimited: keep tokens at +Inf so Allow always returns true.
		b.tokens = math.Inf(1)
	}
	return b
}

func clampInit(initial float64, burst int) float64 {
	if burst <= 0 {
		return 0
	}
	if initial <= 0 {
		return float64(burst)
	}
	if initial > float64(burst) {
		return float64(burst)
	}
	return initial
}

// refill updates tokens based on elapsed time.
func (b *tokenBucket) refill(now time.Time) {
	if b.rate <= 0 || math.IsInf(b.tokens, 1) {
		// unlimited or already infinite
		b.lastRefill = now
		return
	}
	elapsed := now.Sub(b.lastRefill).Seconds()
	if elapsed <= 0 {
		return
	}
	b.tokens = math.Min(b.burst, b.tokens+elapsed*b.rate)
	b.lastRefill = now
}

// hasAtLeast returns whether bucket has >= cost tokens (without consuming).
func (b *tokenBucket) hasAtLeast(now time.Time, cost float64) bool {
	b.refill(now)
	return b.tokens >= cost
}

// consumeNoCheck subtracts cost without re-checking. Caller must ensure hasAtLeast.
func (b *tokenBucket) consumeNoCheck(cost float64) {
	if math.IsInf(b.tokens, 1) {
		return
	}
	b.tokens -= cost
	if b.tokens < 0 {
		b.tokens = 0
	}
}

// take tries to consume 'cost' token(s) atomically.
func (b *tokenBucket) take(now time.Time, cost float64) bool {
	b.refill(now)
	if b.tokens >= cost {
		if !math.IsInf(b.tokens, 1) {
			b.tokens -= cost
		}
		return true
	}
	return false
}
