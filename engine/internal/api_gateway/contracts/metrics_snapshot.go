package contracts

// MetricsSnapshot defines JSON output for /metrics.
type MetricsSnapshot struct {
	RequestsTotal int64   `json:"requests_total"`
	SuccessRate   float64 `json:"success_rate"`
	AvgLatencyMs  float64 `json:"avg_latency_ms"`
	ActiveWorkers int64   `json:"active_workers"`
	QueueDepth    int64   `json:"queue_depth"`
}
