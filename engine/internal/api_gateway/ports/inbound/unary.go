package inbound

import "context"

// Unary allows for unification across mtultiple network transport protocols.
type (
	// UnaryHandler defines a handler for unary requests.
	UnaryHandler[Req any, Res any] func(ctx context.Context, meta RequestMeta, req Req) (Res, error)
	// UnaryMiddleware defines a middleware for unary handlers.
	UnaryMiddleware[Req any, Res any] func(next UnaryHandler[Req, Res]) UnaryHandler[Req, Res]
)
