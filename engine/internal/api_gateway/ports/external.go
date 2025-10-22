package ports

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/types"
)

// Limiter defines rate limiting behaviour.
type Limiter interface {
	Allow(clientID string) bool
}

// Idempotency defines caching for idempotency keys.
type Idempotency interface {
	Get(key string) (types.TransferResult, bool)
	Store(key string, res types.TransferResult)
}

// Dispatcher submits jobs to the worker pool.
type Dispatcher interface {
	Submit(ctx context.Context, cmd types.TransferCommand) <-chan types.TransferResult
	QueueDepth() int64
	ActiveWorkers() int64
}
