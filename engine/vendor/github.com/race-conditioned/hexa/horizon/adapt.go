package horizon

import (
	"github.com/race-conditioned/hexa/apperr"
	"github.com/race-conditioned/hexa/horizon/ports/inbound"
)

// Adapt turns a strongly-typed endpoint into a gateway-compatible erased handler.
func Adapt[c inbound.Ctx, Com inbound.Command, Res inbound.Result](
	h inbound.UnaryHandler[c, Com, Res],
) inbound.UnaryHandler[c, inbound.Command, inbound.Result] {
	return func(ctx c, meta inbound.RequestMeta, req inbound.Command) (inbound.Result, error) {
		typedReq, ok := req.(Com)
		if !ok {
			// Wrong payload wired to this route
			return nil, apperr.Invalid("wrong command type for handler")
		}
		res, err := h(ctx, meta, typedReq)
		if err != nil {
			return nil, err
		}
		// Upcast endpoint-specific result to the erased inbound.Result
		var out inbound.Result = res
		return out, nil
	}
}
