package limiter

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/ports/outbound"
	"time"
)

// Compile-time check that *Composite implements outbound.Limiter.
var _ outbound.Limiter = (*Composite)(nil)

// Composite is a limiter that enforces per-client and optional global limits.
// It is safe for concurrent use.
type Composite struct {
	clock  Clock
	client *clientShardSet
	global *tokenBucket

	cfg             Config
	cleanupTicker   *time.Ticker
	stopCleanupChan chan struct{}
}

// New creates a new Composite limiter. Provide a context that is cancelled on
// server shutdown to stop background cleanup.
func New(ctx context.Context, cfg Config) *Composite {
	clk := cfg.Clock
	if clk == nil {
		clk = systemClock{}
	}
	if cfg.NumShards <= 0 {
		cfg.NumShards = 64
	}
	if cfg.CleanupInterval <= 0 {
		cfg.CleanupInterval = time.Minute
	}
	if cfg.PerClient.TTL <= 0 {
		cfg.PerClient.TTL = 10 * time.Minute
	}

	l := &Composite{
		clock:           clk,
		client:          newClientShardSet(cfg.NumShards, clk, cfg.PerClient),
		cfg:             cfg,
		stopCleanupChan: make(chan struct{}),
	}
	if cfg.Global != nil {
		init := cfg.Global.InitialTokens
		if init <= 0 {
			init = float64(max(1, cfg.Global.Burst))
		}
		l.global = newTokenBucket(clk.Now(), cfg.Global.RatePerSec, cfg.Global.Burst, init)
	}

	// Background janitor: evict idle clients to bound memory.
	l.cleanupTicker = time.NewTicker(cfg.CleanupInterval)
	go func() {
		defer l.cleanupTicker.Stop()
		for {
			select {
			case <-l.cleanupTicker.C:
				l.client.cleanup(cfg.PerClient.TTL)
			case <-ctx.Done():
				return
			case <-l.stopCleanupChan:
				return
			}
		}
	}()
	return l
}

// Stop stops the background cleanup goroutine (optional; otherwise it exits when ctx is cancelled).
func (l *Composite) Stop() { close(l.stopCleanupChan) }

// Allow implements outbound.Limiter.
// Cost is 1 token per call; both global and per-client buckets must admit.
func (l *Composite) Allow(clientID string) bool {
	now := l.clock.Now()

	// Fast path: per-client only.
	if l.global == nil {
		cb := l.client.getOrCreate(clientID)
		cb.b.mu.Lock()
		allowed := cb.b.take(now, 1)
		if allowed {
			cb.lastSeen = now
		}
		cb.b.mu.Unlock()
		return allowed
	}

	// Global + per-client with lock ordering: global.mu -> client.mu
	cb := l.client.getOrCreate(clientID)

	l.global.mu.Lock()
	cb.b.mu.Lock()

	// Lazy refill both; consume only if BOTH have tokens.
	gHas := l.global.hasAtLeast(now, 1)
	cHas := cb.b.hasAtLeast(now, 1)
	if gHas && cHas {
		l.global.consumeNoCheck(1)
		cb.b.consumeNoCheck(1)
		cb.lastSeen = now
		cb.b.mu.Unlock()
		l.global.mu.Unlock()
		return true
	}

	cb.b.mu.Unlock()
	l.global.mu.Unlock()
	return false
}

// Stats returns a snapshot for observability.
// NOTE: Not part of the outbound.Limiter interface; use when you hold the concrete type.
type Stats struct {
	PerClientRate  float64
	PerClientBurst int
	GlobalRate     float64
	GlobalBurst    int
	NumShards      int
}

func (l *Composite) Stats() Stats {
	st := Stats{
		PerClientRate:  l.cfg.PerClient.RatePerSec,
		PerClientBurst: l.cfg.PerClient.Burst,
		NumShards:      l.cfg.NumShards,
	}
	if l.cfg.Global != nil {
		st.GlobalRate = l.cfg.Global.RatePerSec
		st.GlobalBurst = l.cfg.Global.Burst
	}
	return st
}
