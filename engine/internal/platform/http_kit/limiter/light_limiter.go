package limiter

import (
	"context"
	"sync"
	"time"
)

// LightLimiter provides a tiny local rate limiter for basic DoS protection.
type LightLimiter struct {
	rate   float64 // tokens per second
	burst  float64
	mu     sync.Mutex
	bucket map[string]*tokenBucket
}

type tokenBucket struct {
	tokens float64
	last   time.Time
}

// NewLightLimiter returns a local, per-IP limiter (e.g., 10 rps, burst 20).
func NewLightLimiter(rate float64, burst int) *LightLimiter {
	return &LightLimiter{
		rate:   rate,
		burst:  float64(burst),
		bucket: make(map[string]*tokenBucket, 1024),
	}
}

// Allow returns true if the given key is allowed.
func (l *LightLimiter) Allow(_ context.Context, key string) bool {
	now := time.Now()

	l.mu.Lock()
	b, ok := l.bucket[key]
	if !ok {
		b = &tokenBucket{tokens: l.burst, last: now}
		l.bucket[key] = b
	}
	elapsed := now.Sub(b.last).Seconds()
	b.tokens += elapsed * l.rate
	if b.tokens > l.burst {
		b.tokens = l.burst
	}
	b.last = now
	allowed := b.tokens >= 1
	if allowed {
		b.tokens -= 1
	}
	l.mu.Unlock()

	return allowed
}
