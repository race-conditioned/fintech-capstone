package types

import "github.com/google/uuid"

// TransferCommand defines the external API payload for /transfer from any transport.
type TransferCommand struct {
	FromAccount    string
	ToAccount      string
	AmountCents    int64
	IdempotencyKey string
}

// Result returned by the use case.
// It is the internal representation of a completed transfer job.
type TransferResult struct {
	TransactionID uuid.UUID
	Status        string // "success" | "rejected" | "rate_limited" | "duplicate"
	Message       string
}

// TransferResponse is the standard response body for a transfer.
// It is the Result returned from the usecase
type TransferResponse struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	Status        string    `json:"status"` // success | rejected | rate_limited | duplicate
	Message       string    `json:"message,omitempty"`
}

// MetricsSnapshot defines JSON output for /metrics.
type MetricsSnapshot struct {
	RequestsTotal int64   `json:"requests_total"`
	SuccessRate   float64 `json:"success_rate"`
	AvgLatencyMs  float64 `json:"avg_latency_ms"`
	ActiveWorkers int64   `json:"active_workers"`
	QueueDepth    int64   `json:"queue_depth"`
}
