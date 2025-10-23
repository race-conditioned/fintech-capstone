package policy

import (
	"context"
	"errors"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/api_gateway/ports/outbound"
	"fintech-capstone/m/v2/internal/platform/apperr"
	"fmt"
	"time"
)

func RateLimitTransfers(limiter outbound.Limiter, metrics outbound.Metrics) inbound.UnaryMiddleware[inbound.TransferCommand, inbound.TransferResult] {
	return func(next inbound.UnaryHandler[inbound.TransferCommand, inbound.TransferResult]) inbound.UnaryHandler[inbound.TransferCommand, inbound.TransferResult] {
		return func(ctx context.Context, meta inbound.RequestMeta, cmd inbound.TransferCommand) (inbound.TransferResult, error) {
			fmt.Println("Rate Limit")
			if !limiter.Allow(meta.ClientID) {
				if metrics != nil {
					metrics.IncRateLimited()
				}
				return inbound.TransferResult{}, apperr.RateLimited("rate limit exceeded")
			}
			return next(ctx, meta, cmd)
		}
	}
}

func IdempotencyTransfers(store outbound.Idempotency, metrics outbound.Metrics) inbound.UnaryMiddleware[inbound.TransferCommand, inbound.TransferResult] {
	return func(next inbound.UnaryHandler[inbound.TransferCommand, inbound.TransferResult]) inbound.UnaryHandler[inbound.TransferCommand, inbound.TransferResult] {
		return func(ctx context.Context, meta inbound.RequestMeta, cmd inbound.TransferCommand) (inbound.TransferResult, error) {
			fmt.Println("Idempo")
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

func TimeoutTransfers(d time.Duration, metrics outbound.Metrics) inbound.UnaryMiddleware[inbound.TransferCommand, inbound.TransferResult] {
	return func(next inbound.UnaryHandler[inbound.TransferCommand, inbound.TransferResult]) inbound.UnaryHandler[inbound.TransferCommand, inbound.TransferResult] {
		return func(ctx context.Context, meta inbound.RequestMeta, cmd inbound.TransferCommand) (inbound.TransferResult, error) {
			fmt.Println("Timeout")
			// No-op if timeout not configured
			if d <= 0 {
				return next(ctx, meta, cmd)
			}

			cctx, cancel := context.WithTimeout(ctx, d)
			defer cancel()

			done := make(chan struct {
				res inbound.TransferResult
				err error
			}, 1)

			// Run the handler in a separate goroutine
			go func() {
				res, err := next(cctx, meta, cmd)
				done <- struct {
					res inbound.TransferResult
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
					return inbound.TransferResult{}, apperr.Timeout("processing timeout")
				}
				// Context canceled for another reason (e.g., parent canceled)
				return inbound.TransferResult{}, apperr.Internal(cctx.Err().Error())

			case out := <-done:
				// Return successful result if completed before timeout
				return out.res, out.err
			}
		}
	}
}

func ObserveLatencyTransfers(metrics outbound.Metrics) inbound.UnaryMiddleware[inbound.TransferCommand, inbound.TransferResult] {
	return func(next inbound.UnaryHandler[inbound.TransferCommand, inbound.TransferResult]) inbound.UnaryHandler[inbound.TransferCommand, inbound.TransferResult] {
		return func(ctx context.Context, meta inbound.RequestMeta, cmd inbound.TransferCommand) (inbound.TransferResult, error) {
			fmt.Println("Observe Latency")
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
