package inbound

// type Encoder interface {
// 	Encode(writer http.ResponseWriter, status int, v any)
// }

type Writer interface {
	Protocol() string // optional, but handy
	Write(status ResultStatus, v any)
}

// in generic/ports/inbound
type Sink interface {
	// e.g. "http", "grpc" — optional but handy
	Protocol() string
	Write(status string, v any)
}

// Result represents the ubiquitous outcome of processing a Command.
// Any ubiquitous result can be added here.
type Result interface {
	Status() ResultStatus
	Message() string
	Encode(sink Sink)
}

// ResultStatus represents the status of a transfer operation.
type ResultStatus string

const (
	ResultStatusSuccess     ResultStatus = "success"
	ResultStatusRejected    ResultStatus = "rejected"
	ResultStatusRateLimited ResultStatus = "rate_limited"
	ResultStatusDuplicate   ResultStatus = "duplicate"
)

// String returns the string representation of the ResultStatus.
func (ts ResultStatus) String() string {
	return string(ts)
}
