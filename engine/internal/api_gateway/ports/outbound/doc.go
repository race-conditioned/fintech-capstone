// Package outbound declares hexagonal outbound ports the application depends on:
// Dispatcher (worker pool), Limiter (domain rate limit), Idempotency store, and Metrics.
// Concrete adapters live outside this package.
package outbound
