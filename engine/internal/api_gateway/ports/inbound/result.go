package inbound

// Result represents the ubiquitous outcome of processing a Command.
// Any ubiquitous result can be added here.
type Result interface {
	Status() ResultStatus
	Message() string
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
