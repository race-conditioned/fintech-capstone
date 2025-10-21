# Fintech Transaction Processor (Go Concurrency and system design Capstone)

## Motivation

Goal

- Build a concurrent, HTTP-based financial transaction processor in Go that demonstrates real-world concurrency, correctness, and systems scaling.
- Showcase advanced concurrency and reliability engineering through idempotency, rate limiting, autoscaling, profiling, and observability.

Why

- To apply advanced Go concurrency patterns (goroutines, channels, contexts, sync primitives) in a realistic financial domain.
- To showcase understandings of the entire lifecycle of a distributed system: load balancing, backpressure, scaling, fault tolerance, and profiling.

## MVP Definition of Done

### Core Functional Requirements

- A HTTP API (/transfer) that accepts money transfers (from, to, amount) as JSON.
- Each request carries a unique idempotency key; duplicate keys return the original result.
- A ledger service maintains account balances in a thread-safe, atomic manner.
- Transfers are atomic: never overdraw, double-debit, or partial-complete.
- The system can process many (thousands?) of concurrent requests safely.

### Concurrency & System Behavior

- Uses goroutines, channels, and contexts throughout (per request, per worker, cancellation).
- Implements a worker pool that dynamically scales with queue depth.
- Includes backpressure and queue monitoring (no unbounded goroutines).
- Applies per-client and global rate limiting using token buckets. ------------------
- Handles transient errors via exponential backoff + jitter retry (client-side).

### FinTech-Specific Guarantees

- Idempotency: same key → same result (no double-processing).
- Atomic ledger updates using proper locking / lock ordering.
- Precision-safe math (no floats, integer-based cents).
- Consistent lock ordering to avoid deadlocks.
- Ledger integrity validated after concurrent load.

### Observability & Profiling

- /metrics endpoint exposing JSON stats:

```json
{
  "requests_total": 10000,
  "success_rate": 99.5,
  "avg_latency_ms": 120,
  "active_workers": 8,
  "queue_depth": 5
}
```

- /debug/pprof/ endpoint for profiling (CPU, heap, goroutines).
- Flame graph generated from real profiling run (go tool pprof).
- Logging of queue depth, worker utilization, and scaling actions.

### Load Simulation

A concurrent client simulator that:

- Spawns multiple client goroutines.
- Randomly sends transfers between accounts.
- Retries failed transactions with backoff + jitter. (some polite, some not polite)
- Simulates variable request patterns and occasional overload.

## Possibly Future Features

Advanced Systems Concepts (Stretch Goals)

- Sharding / Partitioned Ledger
- Each ledger shard manages a subset of accounts.
- Requests routed by hash(accountID) % N.
- Batch Processing
- Workers can process transactions in small batches for throughput efficiency.
- Circuit Breaker
- Prevents retry storms; temporarily opens when downstream repeatedly fails.
- Bulkheading
- Independent worker pools per account segment to isolate failures.
- Adaptive Feedback Loop
- Adjust rate limits and worker pool sizes based on observed latency and queue depth.
- Event Sourcing / Write-Ahead Log
- Ledger writes appended to immutable event log for replay/recovery.
- Replay Safety
- On restart, replay transactions safely (idempotent by key).
- p99 Latency Tracking
- Record 50th, 90th, 99th percentile latencies for realistic service-level objectives.
- Tracing & Context Propagation
- Assign a trace ID per request, propagated through all goroutines.
- Persistent Storage Integration
- Swap in SQLite / BoltDB backend for persistence.
- Security Layer
- API key-based authentication, HMAC request signing, and replay protection.
- Chaos Mode
- Random latency spikes, dropped requests, and simulated component failures.

## Knowledge Gaps

To master during or after the capstone:

- Advanced profiling – interpreting flame graphs, detecting contention and blocking.
- Atomic operations – sync/atomic and lock-free counters.
- Lock ordering and deadlock prevention patterns.
- Circuit breaker and bulkhead patterns in Go.
- Prometheus / expvar integration for real metrics.
- Dynamic rate limiting (PID-style feedback loop tuning).
- OpenTelemetry tracing: distributed trace correlation.
- Consistency models: eventual vs strong consistency trade-offs.
- (metanote: update this list as I identify new system-level learnings during implementation)

## System Phases / Roadmap

Phase Focus Outcome

1. API & Ledger Build /transfer endpoint, validate input, update ledger atomically Working, correct base system
2. Worker Pool & Queues Introduce async job processing & dynamic autoscaling Stable concurrency core
3. Idempotency & Atomicity Ensure safe replay & exactly-once semantics No double-spend or race conditions
4. Rate Limiting Per-client & global token buckets Fair & resilient load control
5. Client Simulator Generate concurrent load w/ retries & jitter Stress test concurrency correctness
6. Profiling & Metrics Add pprof, expvar, structured logs Measure and visualize performance
7. Reliability Features Circuit breaker, chaos, batching, sharding Production-grade resilience
8. Optimization & Tuning Analyze flame graphs, reduce contention Measured scalability and insight
   Deliverables

## Initial directory structure sketch

/cmd/server/main.go - + cobra commands
/cmd/admin/main.go - + cobra commands
/internal/store/
/internal/ledger/
/internal/api_gateway/
/internal/worker_pool/
/internal/limiter/
/internal/idempotency/
/internal/metrics/
/internal/clientload/
/internal/cache/
/internal/admin/
/internal/metrics/
/internal/persistence/

## Deliverables

- README with architecture diagram & profiling examples.
- pprof flame graph screenshots (Stretch: Visual pprof flame graph (go tool pprof -http=:8081 cpu.prof))
- Metrics output under load test.
- Optional: CSV log for latency analysis.
- Comparison of throughput pre- and post-optimization.
- Load testing summary:

```json
Load test: 100 clients, 10,000 txns
Success: 99.4%
RateLimited: 3.2%
Avg Latency: 128ms
p99 Latency: 280ms
Max Workers: 8
```
