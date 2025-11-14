package horizon

import (
	"sync"

	"github.com/race-conditioned/hexa/horizon/ports/inbound"
)

type HandlerKey string

func (k HandlerKey) String() string {
	return string(k)
}

// Option configures a Gateway.
type Option[ctx inbound.Ctx] func(*Gateway[ctx])

// Gateway is the API Gateway entrypoint, composing handlers with middleware.
type Gateway[ctx inbound.Ctx] struct {
	mu         sync.RWMutex
	handlers   map[HandlerKey]inbound.UnaryHandler[ctx, inbound.Command, inbound.Result]
	middleware []inbound.UnaryMiddleware[ctx, inbound.Command, inbound.Result]
	plugins    ctx
}

// NewGateway constructs a new API Gateway entrypoint with all handlers composed with middleware.
func NewGateway[c inbound.Ctx](
	ctx c,
) *Gateway[c] {
	return &Gateway[c]{
		handlers:   make(map[HandlerKey]inbound.UnaryHandler[c, inbound.Command, inbound.Result]),
		middleware: make([]inbound.UnaryMiddleware[c, inbound.Command, inbound.Result], 0),
		plugins:    ctx,
	}
}

func (g *Gateway[ctx]) RegisterHandler(name HandlerKey, h inbound.UnaryHandler[ctx, inbound.Command, inbound.Result]) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.handlers[name] = h
}

func (g *Gateway[ctx]) Handler(name HandlerKey) (inbound.UnaryHandler[ctx, inbound.Command, inbound.Result], bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	h, ok := g.handlers[name]
	return h, ok
}

func (g *Gateway[ctx]) Handlers() map[HandlerKey]inbound.UnaryHandler[ctx, inbound.Command, inbound.Result] {
	handlerMap := make(map[HandlerKey]inbound.UnaryHandler[ctx, inbound.Command, inbound.Result])
	g.mu.RLock()
	defer g.mu.RUnlock()
	for k, v := range g.handlers {
		handlerMap[k] = v
	}
	return handlerMap
}
