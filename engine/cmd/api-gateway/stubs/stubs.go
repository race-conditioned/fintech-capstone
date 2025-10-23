package stubs

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/contracts"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/api_gateway/ports/outbound"
	"sync"
	"time"

	"github.com/google/uuid"
)

func Build() (outbound.Limiter, outbound.Idempotency, outbound.Dispatcher, outbound.Metrics) {
	return &allowAllLimiter{}, newInmemIdemp(), &immediateDispatcher{}, &noopMetrics{}
}

// Limiter: allow all
type allowAllLimiter struct{}

func (a *allowAllLimiter) Allow(string) bool { return true }

// Idempotency: in-memory
type inmemIdemp struct {
	mu sync.Mutex
	m  map[string]inbound.TransferResult
}

func newInmemIdemp() *inmemIdemp { return &inmemIdemp{m: make(map[string]inbound.TransferResult)} }
func (s *inmemIdemp) Get(k string) (inbound.TransferResult, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.m[k]
	return v, ok
}

func (s *inmemIdemp) Store(k string, v inbound.TransferResult) {
	s.mu.Lock()
	s.m[k] = v
	s.mu.Unlock()
}

// Dispatcher: immediately succeed
type immediateDispatcher struct{}

func (d *immediateDispatcher) Submit(_ context.Context, _ inbound.TransferCommand) <-chan inbound.TransferResult {
	ch := make(chan inbound.TransferResult, 1)
	ch <- inbound.TransferResult{TransactionID: uuid.New(), Status: "success", Message: "ok"}
	close(ch)
	return ch
}
func (d *immediateDispatcher) QueueDepth() int64    { return 0 }
func (d *immediateDispatcher) ActiveWorkers() int64 { return 0 }

// Metrics: no-op
type noopMetrics struct{}

func (*noopMetrics) IncRequest()                         {}
func (*noopMetrics) IncSuccess()                         {}
func (*noopMetrics) IncRateLimited()                     {}
func (*noopMetrics) IncTimeout()                         {}
func (*noopMetrics) IncIdempotentHit()                   {}
func (*noopMetrics) ObserveLatency(time.Duration)        {}
func (*noopMetrics) Snapshot() contracts.MetricsSnapshot { return contracts.MetricsSnapshot{} }
