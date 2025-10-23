package inbound

import "context"

// Server defines the interface for inbound network transport servers.
type Server interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
