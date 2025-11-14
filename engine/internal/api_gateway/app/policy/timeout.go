package policy

import (
	"context"
	"errors"
	"fmt"

	"fintech-capstone/m/v2/internal/platform/apperr"

	"github.com/race-conditioned/hexa/horizon/ports/inbound"
)

// Timeout is a middleware that enforces a timeout on request processing.
func Timeout(next AppHandler) AppHandler {
	return func(ctx Plugins, meta inbound.RequestMeta, cmd inbound.Command) (inbound.Result, error) {
		fmt.Println("Applying timeout...")
		var zero inbound.Result
		// No-op if timeout not configured
		if ctx.Timeout() <= 0 {
			return next(ctx, meta, cmd)
		}
		cctx, cancel := context.WithTimeout(context.Background(), ctx.Timeout()) // TODO: propogate context in horizon.context
		defer cancel()

		done := make(chan struct {
			res inbound.Result
			err error
		}, 1)

		go func() {
			r, e := next(ctx, meta, cmd) // TODO: add new ctx, with timeout
			done <- struct {
				res inbound.Result
				err error
			}{r, e}
		}()

		select {
		case <-cctx.Done():
			// Explicit timeout: propagate a well-defined app error
			if errors.Is(cctx.Err(), context.DeadlineExceeded) {
				ctx.Metrics().IncTimeout()
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
