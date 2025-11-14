package inbound

import "context"

type Ctx interface {
	context.Context
}

// Unary allows for unification across multiple network transport protocols.
type (
	// UnaryHandler defines a handler for unary requests.
	UnaryHandler[c Ctx, com Command, res Result] func(ctx c, meta RequestMeta, req com) (res, error)
	// UnaryMiddleware defines a middleware for unary handlers.
	UnaryMiddleware[c Ctx, com Command, res Result] func(next UnaryHandler[c, com, res]) UnaryHandler[c, com, res]
)
