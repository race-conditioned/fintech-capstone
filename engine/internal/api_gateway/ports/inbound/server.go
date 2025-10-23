package inbound

import "context"

type Server interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
