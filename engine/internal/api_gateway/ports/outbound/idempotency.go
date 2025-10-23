package outbound

// Idempotency defines caching for idempotency keys.
type Idempotency[T any] interface {
	Get(key string) (T, bool)
	Store(key string, res T)
}
