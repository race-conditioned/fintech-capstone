package ports

import "context"

type InboundServer interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
