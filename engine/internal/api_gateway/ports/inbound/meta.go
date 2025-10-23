package inbound

type RequestMeta struct {
	ClientID  string
	RequestID string
	TraceID   string
	RemoteIP  string
	Protocol  string // "http","grpc",...
	Target    string // path or method name for logging
}
