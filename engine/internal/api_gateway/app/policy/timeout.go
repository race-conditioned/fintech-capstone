package policy

import (
	"context"
	"errors"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/api_gateway/ports/outbound"
	"fintech-capstone/m/v2/internal/platform/apperr"
	"time"
)

// Timeout is a middleware that enforces a timeout on request processing.
func Timeout[Com inbound.Command, Res inbound.Result](d time.Duration, counterMetrics outbound.CounterMetrics) inbound.UnaryMiddleware[Com, Res] {
	return func(next inbound.UnaryHandler[Com, Res]) inbound.UnaryHandler[Com, Res] {
		return func(ctx context.Context, meta inbound.RequestMeta, cmd Com) (Res, error) {
			var zero Res
			// No-op if timeout not configured
			if d <= 0 {
				return next(ctx, meta, cmd)
			}
			cctx, cancel := context.WithTimeout(ctx, d)
			defer cancel()

			done := make(chan struct {
				res Res
				err error
			}, 1)

			go func() {
				r, e := next(cctx, meta, cmd)
				done <- struct {
					res Res
					err error
				}{r, e}
			}()

			select {
			case <-cctx.Done():
				// Explicit timeout: propagate a well-defined app error
				if errors.Is(cctx.Err(), context.DeadlineExceeded) {
					if counterMetrics != nil {
						counterMetrics.IncTimeout()
					}
					return zero, apperr.Timeout("processing timeout")
				}
				// Context canceled for another reason (e.g., parent canceled)
				return zero, apperr.Internal(cctx.Err().Error())
			case out := <-done:
				// Return successful result if completed before timeout
				return out.res, out.err
			}
		}
	}
}
