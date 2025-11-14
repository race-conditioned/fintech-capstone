package inbound

// Command is a base interface for all request commands.
// Any ubiquitous capability can be added here.
type Command interface{}

type CommandDTO interface {
	ToCommand() Command
}

type RequestMeta struct {
	ClientID  string
	RequestID string
	TraceID   string
	RemoteIP  string
	Protocol  string
	Target    string
}
