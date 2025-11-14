// Package intake defines the generic ingress policy spec.
// It is interpreted by specific transport adapters (HTTP, gRPC, GraphQL, etc.)
package intake

import "context"

// Spec describes ingress-level policies such as
// request size, concurrency, panic recovery, and logging.
// It contains *transport-neutral* settings and hook slots.
// Adapters map them to their native mechanisms.
type Spec struct {
	// --- Wire-level guards ---
	MaxBodyBytes       int64 // 0 = unlimited
	MaxInFlight        int   // 0 = unlimited
	EnableRecover      bool
	EnableReqID        bool
	EnableReqLog       bool
	EnableRateLimiting bool

	// --- Behavioural hooks (optional) ---
	OnBodyTooLarge    Hook // e.g. called when payload exceeds limit
	OnTooManyInFlight Hook
	OnRecover         RecoverHook   // rec is the recovered value
	OnRequestID       RequestIDHook // when a new or propagated ID is known
	OnLog             LogHook       // per-request log summary
	OnRateLimit       Hook

	// --- Rate limiting function (optional) ---
	RateLimiter AllowFunc // nil = skip

	// --- Middleware chain extension (adapter-specific) ---
	Middleware []MWHook // adapters interpret element types they understand

	// --- App-level defaults ---
	DefaultTimeoutMs int
}

// AllowFunc is a minimal adapter to allow various rate limiting decisions.
type AllowFunc func(ctx context.Context, key string) bool

// Hook is a generic callback for simple events.
type Hook func(Event)

type Timing string

const (
	PreHook  Timing = "pre"
	PostHook Timing = "post"
)

type MWHook struct {
	Fn     func(Event)
	Timing Timing
}

// RecoverHook is called when a panic is recovered.
type RecoverHook func(Event, any)

// RequestIDHook reports a generated or extracted request ID.
type RequestIDHook func(Event, string)

// LogHook provides structured logging metadata.
type LogHook func(Event, LogFields)

// Event is a small, generic metadata bag shared across transports.
type Event struct {
	Target   string // route, method, RPC name, etc.
	ClientID string // extracted client identity if known
	Protocol string // "http", "grpc", "graphql", etc.
}

// LogFields carries neutral log data.
type LogFields struct {
	StatusCode int
	Bytes      int64
	LatencyMs  int64
	Method     string
	RequestID  string
	RemoteAddr string
	UserAgent  string
}
