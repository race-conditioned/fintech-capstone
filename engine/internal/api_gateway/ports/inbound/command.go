package inbound

// Command is a base interface for all request commands.
// Any ubiquitous capability can be added here.
type Command interface{}

// Idempotent is an optional Command Capability
type Idempotent interface {
	IdempotencyKey() string
}

// IdempotentCommand is a Command that supports Idempotency
type IdempotentCommand interface {
	Command
	Idempotent
}
