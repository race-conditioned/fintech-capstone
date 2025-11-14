package symphony

import (
	"sync"

	"github.com/race-conditioned/hexa/horizon/ports/inbound"
)

// SymphonyOption customizes the Symphony composer (pre/post only, universal).
type SymphonyOption[c inbound.Ctx] func(*Symphony[c])

// Symphony holds global (type-erased) pre/post middleware and their orders.
type Symphony[c inbound.Ctx] struct {
	mu        sync.RWMutex
	preOrder  PolicyOrder
	postOrder PolicyOrder
	pre       map[PolicyStage]inbound.UnaryMiddleware[c, inbound.Command, inbound.Result]
	post      map[PolicyStage]inbound.UnaryMiddleware[c, inbound.Command, inbound.Result]
}

// New creates a Symphony with explicit Pre/Post orders.
func New[c inbound.Ctx](
	preOrder PolicyOrder,
	postOrder PolicyOrder,
	opts ...SymphonyOption[c],
) *Symphony[c] {
	s := &Symphony[c]{
		preOrder:  preOrder,
		postOrder: postOrder,
		pre:       make(map[PolicyStage]inbound.UnaryMiddleware[c, inbound.Command, inbound.Result]),
		post:      make(map[PolicyStage]inbound.UnaryMiddleware[c, inbound.Command, inbound.Result]),
	}

	for _, opt := range opts {
		opt(s)
	}
	return s
}

// WithPre registers a global pre middleware by stage.
func WithPre[c inbound.Ctx](stage PolicyStage, mw inbound.UnaryMiddleware[c, inbound.Command, inbound.Result]) SymphonyOption[c] {
	return func(s *Symphony[c]) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.pre[stage] = mw
	}
}

// WithPost registers a global post middleware by stage.
func WithPost[c inbound.Ctx](stage PolicyStage, mw inbound.UnaryMiddleware[c, inbound.Command, inbound.Result]) SymphonyOption[c] {
	return func(s *Symphony[c]) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.post[stage] = mw
	}
}

// -------------------- Composition (per-endpoint pipeline) --------------------

// CompositionOption customizes a per-endpoint composition.
type CompositionOption[c inbound.Ctx, com inbound.Command, res inbound.Result] func(*Composition[c, com, res])

// Composition represents a per-endpoint stack (its own mid-stages, plus shared pre/post).
type Composition[c inbound.Ctx, com inbound.Command, res inbound.Result] struct {
	mu    sync.RWMutex
	order PolicyOrder // mid/policy order, specific to this endpoint
	mws   map[PolicyStage]inbound.UnaryMiddleware[c, com, res]

	// inherited from Symphony (global)
	preOrder  PolicyOrder
	postOrder PolicyOrder
	pre       map[PolicyStage]inbound.UnaryMiddleware[c, inbound.Command, inbound.Result]
	post      map[PolicyStage]inbound.UnaryMiddleware[c, inbound.Command, inbound.Result]
}

// Compose creates a per-endpoint composition with its own middle-stage order.
// (Top-level generic function because methods can’t have their own type params.)
func Compose[c inbound.Ctx, com inbound.Command, res inbound.Result](
	s *Symphony[c],
	order PolicyOrder,
	opts ...CompositionOption[c, com, res],
) *Composition[c, com, res] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	comp := &Composition[c, com, res]{
		order:     order,
		mws:       make(map[PolicyStage]inbound.UnaryMiddleware[c, com, res]),
		preOrder:  s.preOrder,
		postOrder: s.postOrder,
		pre:       s.pre,
		post:      s.post,
	}
	for _, opt := range opts {
		opt(comp)
	}
	return comp
}

// WithPolicy registers a middleware for a specific mid-stage on this composition.
func WithPolicy[c inbound.Ctx, com inbound.Command, res inbound.Result](
	stage PolicyStage,
	mw inbound.UnaryMiddleware[c, com, res],
) CompositionOption[c, com, res] {
	return func(c *Composition[c, com, res]) {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.mws[stage] = mw
	}
}

// Wrap applies: post (outermost-after), then mid (ordered), then pre (outermost-before).
func (c *Composition[ctx, com, res]) Wrap(base inbound.UnaryHandler[ctx, com, res]) inbound.UnaryHandler[ctx, com, res] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	h := base

	// 1) Post
	if c.postOrder != nil {
		ord := c.postOrder.Order()
		for i := len(ord) - 1; i >= 0; i-- {
			if mwU, ok := c.post[ord[i]]; ok {
				h = Lift[ctx, com, res](mwU)(h)
			}
		}
	}

	// 2) Mid/policies
	if c.order != nil {
		ord := c.order.Order()
		for i := len(ord) - 1; i >= 0; i-- {
			if mw, ok := c.mws[ord[i]]; ok {
				h = mw(h)
			}
		}
	}

	// 3) Pre
	if c.preOrder != nil {
		ord := c.preOrder.Order()
		for i := len(ord) - 1; i >= 0; i-- {
			if mwU, ok := c.pre[ord[i]]; ok {
				h = Lift[ctx, com, res](mwU)(h)
			}
		}
	}
	return h
}

// Chain provides a simple one-off manual chain utility.
func Chain[c inbound.Ctx](
	base inbound.UnaryHandler[c, inbound.Command, inbound.Result],
	mws ...inbound.UnaryMiddleware[c, inbound.Command, inbound.Result],
) inbound.UnaryHandler[c, inbound.Command, inbound.Result] {
	h := base
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
