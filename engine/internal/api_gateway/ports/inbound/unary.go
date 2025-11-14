package inbound

// // Unary allows for unification across multiple network transport protocols.
// type (
// 	// UnaryHandler defines a handler for unary requests.
// 	UnaryHandler[Req Command, Res Result] func(ctx context.Context, meta RequestMeta, req Req) (Res, error)
// 	// UnaryMiddleware defines a middleware for unary handlers.
// 	UnaryMiddleware[Req Command, Res Result] func(next UnaryHandler[Req, Res]) UnaryHandler[Req, Res]
// )
