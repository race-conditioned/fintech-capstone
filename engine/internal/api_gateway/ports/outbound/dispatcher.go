package outbound

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
)

// Dispatcher submits jobs to the worker pool.
type Dispatcher interface {
	Submit(ctx context.Context, cmd inbound.TransferCommand) inbound.TransferResult
	QueueDepth() int64
	ActiveWorkers() int64
}
