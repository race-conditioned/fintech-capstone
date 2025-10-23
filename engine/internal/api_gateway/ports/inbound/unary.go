package inbound

import "context"

type (
	UnaryHandler[Req any, Res any]    func(ctx context.Context, meta RequestMeta, req Req) (Res, error)
	UnaryMiddleware[Req any, Res any] func(next UnaryHandler[Req, Res]) UnaryHandler[Req, Res]
)
