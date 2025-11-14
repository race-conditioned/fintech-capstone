package outbound

// Limiter defines rate limiting behavior.
type Limiter interface {
	Allow(clientID string) bool
}
