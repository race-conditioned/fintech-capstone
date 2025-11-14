package policy

import (
	"fmt"

	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"

	hexa_inbound "github.com/race-conditioned/hexa/horizon/ports/inbound"
)

type IdempotentHandler = hexa_inbound.UnaryHandler[Plugins, inbound.IdempotentCommand, hexa_inbound.Result]

// Idempotency is a middlelare that provides idempotency support for commands implementing IdempotentCommand.
func Idempotency(next IdempotentHandler) IdempotentHandler {
	return func(ctx Plugins, meta hexa_inbound.RequestMeta, cmd inbound.IdempotentCommand) (hexa_inbound.Result, error) {
		fmt.Println("idempotency")
		// WARN: can store be nil?
		if cached, ok := ctx.Idempotency().Get(cmd.IdempotencyKey()); ok {
			ctx.Metrics().IncIdempotentHit()
			return cached, nil
		}

		res, err := next(ctx, meta, cmd)
		// WARN: haven't checked if idempotencykey can be empty
		if err == nil {
			ctx.Idempotency().Store(cmd.IdempotencyKey(), res)
		}

		return res, err
	}
}
