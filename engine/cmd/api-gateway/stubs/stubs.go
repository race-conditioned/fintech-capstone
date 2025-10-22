package stubs

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/ports"
	"fintech-capstone/m/v2/internal/api_gateway/types"
	"sync"
	"time"

	"github.com/google/uuid"
)

func Build() (ports.Limiter, ports.Idempotency, ports.Dispatcher, ports.Metrics) {
	return &allowAllLimiter{}, newInmemIdemp(), &immediateDispatcher{}, &noopMetrics{}
}

// Limiter: allow all
type allowAllLimiter struct{}

func (a *allowAllLimiter) Allow(string) bool { return true }

// Idempotency: in-memory
type inmemIdemp struct {
	mu sync.Mutex
	m  map[string]types.TransferResult
}

func newInmemIdemp() *inmemIdemp { return &inmemIdemp{m: make(map[string]types.TransferResult)} }
func (s *inmemIdemp) Get(k string) (types.TransferResult, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.m[k]
	return v, ok
}
func (s *inmemIdemp) Store(k string, v types.TransferResult) { s.mu.Lock(); s.m[k] = v; s.mu.Unlock() }

// Dispatcher: immediately succeed
type immediateDispatcher struct{}

func (d *immediateDispatcher) Submit(_ context.Context, _ types.TransferCommand) <-chan types.TransferResult {
	ch := make(chan types.TransferResult, 1)
	ch <- types.TransferResult{TransactionID: uuid.New(), Status: "success", Message: "ok"}
	close(ch)
	return ch
}
func (d *immediateDispatcher) QueueDepth() int64    { return 0 }
func (d *immediateDispatcher) ActiveWorkers() int64 { return 0 }

// Metrics: no-op
type noopMetrics struct{}

func (*noopMetrics) IncRequest()                     {}
func (*noopMetrics) IncSuccess()                     {}
func (*noopMetrics) IncRateLimited()                 {}
func (*noopMetrics) IncTimeout()                     {}
func (*noopMetrics) IncIdempotentHit()               {}
func (*noopMetrics) ObserveLatency(time.Duration)    {}
func (*noopMetrics) Snapshot() types.MetricsSnapshot { return types.MetricsSnapshot{} }
