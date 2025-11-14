package symphony

import (
	"github.com/race-conditioned/hexa/apperr"
	"github.com/race-conditioned/hexa/horizon/ports/inbound"
)

// Lift adapts a universal (type-erased) middleware to a typed one.
func Lift[c inbound.Ctx, Com inbound.Command, Res inbound.Result](
	mw inbound.UnaryMiddleware[c, inbound.Command, inbound.Result],
) inbound.UnaryMiddleware[c, Com, Res] {
	return func(next inbound.UnaryHandler[c, Com, Res]) inbound.UnaryHandler[c, Com, Res] {
		// Erased version of next
		erasedNext := func(ctx c, meta inbound.RequestMeta, req inbound.Command) (inbound.Result, error) {
			typedReq, ok := req.(Com)
			if !ok {
				return nil, apperr.Invalid("wrong command type for middleware")
			}
			res, err := next(ctx, meta, typedReq)
			if err != nil {
				return nil, err
			}
			return res, nil // upcast
		}

		erasedWrapped := mw(erasedNext)

		// Return typed wrapper calling through the erased chain
		return func(ctx c, meta inbound.RequestMeta, cmd Com) (Res, error) {
			r, err := erasedWrapped(ctx, meta, cmd) // cmd upcasts to inbound.Command
			if err != nil {
				var zero Res
				return zero, err
			}
			typedRes, ok := r.(Res)
			if !ok {
				var zero Res
				return zero, apperr.Internal("middleware result type mismatch")
			}
			return typedRes, nil
		}
	}
}

// LiftCap adapts a capability-specific middleware (Cap) that returns
// inbound.Result (erased). If Com doesn't implement Cap, it's a no-op.
//
// Use this for capability middlewares like Idempotency that are written
// against a capability interface and return inbound.Result.
func LiftCap[c inbound.Ctx, Com inbound.Command, Res inbound.Result, Cap inbound.Command](
	mw inbound.UnaryMiddleware[c, Cap, inbound.Result],
) inbound.UnaryMiddleware[c, Com, Res] {
	return func(next inbound.UnaryHandler[c, Com, Res]) inbound.UnaryHandler[c, Com, Res] {
		// Bridge next (typed) -> capability next (returns inbound.Result)
		capNext := func(ctx c, meta inbound.RequestMeta, capCmd Cap) (inbound.Result, error) {
			typedCmd, ok := any(capCmd).(Com)
			if !ok {
				return nil, apperr.Internal("capability bridge type mismatch")
			}
			res, err := next(ctx, meta, typedCmd)
			if err != nil {
				return nil, err
			}
			return res, nil // upcast Res -> inbound.Result
		}

		wrappedCap := mw(capNext)

		return func(ctx c, meta inbound.RequestMeta, cmd Com) (Res, error) {
			// Only run if the command actually supports the capability.
			if capCmd, ok := any(cmd).(Cap); ok {
				r, err := wrappedCap(ctx, meta, capCmd)
				if err != nil {
					var zero Res
					return zero, err
				}
				typedRes, ok := r.(Res)
				if !ok {
					var zero Res
					return zero, apperr.Internal("capability result type mismatch")
				}
				return typedRes, nil
			}
			// No capability -> pass through
			return next(ctx, meta, cmd)
		}
	}
}
