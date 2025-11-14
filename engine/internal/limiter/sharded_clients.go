package limiter

import (
	"hash/fnv"
	"sync"
	"time"
)

type clientBucket struct {
	b *tokenBucket

	// updated on every Allow() call (shard lock not required to read/write);
	// guarded by the bucket's own mutex for consistency with tokens/lastRefill.
	lastSeen time.Time
}

type shard struct {
	mu   sync.RWMutex
	data map[string]*clientBucket
}

type clientShardSet struct {
	shards []shard
	clock  Clock
	cfg    PerClientConfig
}

func newClientShardSet(n int, clock Clock, cfg PerClientConfig) *clientShardSet {
	if n <= 0 {
		n = 64
	}
	s := make([]shard, n)
	for i := range s {
		s[i].data = make(map[string]*clientBucket, 256)
	}
	return &clientShardSet{
		shards: s,
		clock:  clock,
		cfg:    cfg,
	}
}

func (cs *clientShardSet) getShard(key string) *shard {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	idx := int(h.Sum32()) & (len(cs.shards) - 1) // len is power of 2? If not, mod:
	if (len(cs.shards) & (len(cs.shards) - 1)) != 0 {
		idx = int(h.Sum32()) % len(cs.shards)
	}
	return &cs.shards[idx]
}

// getOrCreate returns the client's bucket (creating if missing).
func (cs *clientShardSet) getOrCreate(clientID string) *clientBucket {
	now := cs.clock.Now()
	sh := cs.getShard(clientID)

	// Fast path read
	sh.mu.RLock()
	cb := sh.data[clientID]
	sh.mu.RUnlock()
	if cb != nil {
		return cb
	}

	// Create under write lock
	sh.mu.Lock()
	cb = sh.data[clientID]
	if cb == nil {
		init := cs.cfg.InitialTokens
		if init <= 0 {
			init = float64(max(1, cs.cfg.Burst))
		}
		cb = &clientBucket{
			b:        newTokenBucket(now, cs.cfg.RatePerSec, cs.cfg.Burst, init),
			lastSeen: now,
		}
		sh.data[clientID] = cb
	}
	sh.mu.Unlock()
	return cb
}

// cleanup evicts idle client buckets across shards.
func (cs *clientShardSet) cleanup(ttl time.Duration) (evicted int) {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	now := cs.clock.Now()
	for i := range cs.shards {
		sh := &cs.shards[i]
		sh.mu.Lock()
		for id, cb := range sh.data {
			cb.b.mu.Lock()
			ls := cb.lastSeen
			cb.b.mu.Unlock()
			if now.Sub(ls) >= ttl {
				delete(sh.data, id)
				evicted++
			}
		}
		sh.mu.Unlock()
	}
	return
}
