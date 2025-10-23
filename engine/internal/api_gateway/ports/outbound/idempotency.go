package outbound

import (
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
)

// Idempotency defines caching for idempotency keys.
type Idempotency interface {
	Get(key string) (inbound.TransferResult, bool)
	Store(key string, res inbound.TransferResult)
}
