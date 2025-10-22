package policy

import (
	"context"
	"errors"
	"fintech-capstone/m/v2/internal/api_gateway/apperr"
	"fintech-capstone/m/v2/internal/api_gateway/ports"
	"fintech-capstone/m/v2/internal/api_gateway/types"
	"time"
)

func RateLimitTransfers(limiter ports.Limiter, metrics ports.Metrics) ports.UnaryMiddleware[types.TransferCommand, types.TransferResult] {
	return func(next ports.UnaryHandler[types.TransferCommand, types.TransferResult]) ports.UnaryHandler[types.TransferCommand, types.TransferResult] {
		return func(ctx context.Context, meta ports.RequestMeta, cmd types.TransferCommand) (types.TransferResult, error) {
			if !limiter.Allow(meta.ClientID) {
				if metrics != nil {
					metrics.IncRateLimited()
				}
				return types.TransferResult{}, apperr.RateLimited("rate limit exceeded")
			}
			return next(ctx, meta, cmd)
		}
	}
}

func IdempotencyTransfers(store ports.Idempotency, metrics ports.Metrics) ports.UnaryMiddleware[types.TransferCommand, types.TransferResult] {
	return func(next ports.UnaryHandler[types.TransferCommand, types.TransferResult]) ports.UnaryHandler[types.TransferCommand, types.TransferResult] {
		return func(ctx context.Context, meta ports.RequestMeta, cmd types.TransferCommand) (types.TransferResult, error) {
			if store != nil {
				if cached, ok := store.Get(cmd.IdempotencyKey); ok {
					if metrics != nil {
						metrics.IncIdempotentHit()
					}
					return cached, nil
				}
			}
			res, err := next(ctx, meta, cmd)
			if err == nil && store != nil && cmd.IdempotencyKey != "" {
				store.Store(cmd.IdempotencyKey, res)
			}
			return res, err
		}
	}
}

func TimeoutTransfers(d time.Duration, metrics ports.Metrics) ports.UnaryMiddleware[types.TransferCommand, types.TransferResult] {
	return func(next ports.UnaryHandler[types.TransferCommand, types.TransferResult]) ports.UnaryHandler[types.TransferCommand, types.TransferResult] {
		return func(ctx context.Context, meta ports.RequestMeta, cmd types.TransferCommand) (types.TransferResult, error) {
			// No-op if timeout not configured
			if d <= 0 {
				return next(ctx, meta, cmd)
			}

			cctx, cancel := context.WithTimeout(ctx, d)
			defer cancel()

			done := make(chan struct {
				res types.TransferResult
				err error
			}, 1)

			// Run the handler in a separate goroutine
			go func() {
				res, err := next(cctx, meta, cmd)
				done <- struct {
					res types.TransferResult
					err error
				}{res, err}
			}()

			select {
			case <-cctx.Done():
				// Explicit timeout: propagate a well-defined app error
				if errors.Is(cctx.Err(), context.DeadlineExceeded) {
					if metrics != nil {
						metrics.IncTimeout()
					}
					return types.TransferResult{}, apperr.Timeout("processing timeout")
				}
				// Context canceled for another reason (e.g., parent canceled)
				return types.TransferResult{}, apperr.Internal(cctx.Err().Error())

			case out := <-done:
				// Return successful result if completed before timeout
				return out.res, out.err
			}
		}
	}
}

func ObserveLatencyTransfers(metrics ports.Metrics) ports.UnaryMiddleware[types.TransferCommand, types.TransferResult] {
	return func(next ports.UnaryHandler[types.TransferCommand, types.TransferResult]) ports.UnaryHandler[types.TransferCommand, types.TransferResult] {
		return func(ctx context.Context, meta ports.RequestMeta, cmd types.TransferCommand) (types.TransferResult, error) {
			start := time.Now()
			res, err := next(ctx, meta, cmd)
			if metrics != nil {
				metrics.ObserveLatency(time.Since(start))
				if err == nil && res.Status == "success" {
					metrics.IncSuccess()
				}
			}
			return res, err
		}
	}
}
