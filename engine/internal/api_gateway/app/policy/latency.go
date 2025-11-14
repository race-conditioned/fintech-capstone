package policy

import (
	"fmt"
	"time"

	"github.com/race-conditioned/hexa/horizon/ports/inbound"
)

// ObserveLatency is a middleware that observes the latency of requests and records success metrics.
func ObserveLatency(next AppHandler) AppHandler {
	return func(ctx Plugins, meta inbound.RequestMeta, cmd inbound.Command) (inbound.Result, error) {
		fmt.Println("Observing latency...")
		start := time.Now()
		res, err := next(ctx, meta, cmd)
		ctx.Metrics().ObserveLatency(time.Since(start))
		return res, err
	}
}
