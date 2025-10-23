package inbound

// NetComProtocol defines supported network communication protocols.
type NetComProtocol string

const (
	ProtocolHTTP NetComProtocol = "http"
	ProtocolGRPC NetComProtocol = "grpc"
)

// RequestMeta holds metadata about an inbound request.
type RequestMeta struct {
	ClientID  string
	RequestID string
	TraceID   string
	RemoteIP  string
	Protocol  NetComProtocol
	Target    string // path or method name for logging
}
