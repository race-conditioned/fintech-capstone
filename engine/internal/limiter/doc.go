// Package limiter provides a high-quality, concurrency-safe token bucket rate
// limiter with per-client and optional global limits. It implements the
// outbound.Limiter interface used by the policy middleware.
//
// Design goals:
//   - Non-blocking (boolean Allow) with deterministic, lazy refill.
//   - Per-client isolation and fairness via sharded maps.
//   - Optional global bucket to cap aggregate QPS.
//   - Monotonic time, burst support, safe under high contention.
//   - Low heap churn & passive cleanup of inactive clients.
//
// Lock ordering (to prevent deadlocks):
//  1. Global bucket mutex (if configured)
//  2. Per-client bucket mutex
//
// Shards are only locked briefly when retrieving/creating client buckets; the
// actual rate decision holds only bucket-level locks.
package limiter
