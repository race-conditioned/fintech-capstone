package outbound

// Limiter defines rate limiting behaviour.
type Limiter interface {
	Allow(clientID string) bool
}
