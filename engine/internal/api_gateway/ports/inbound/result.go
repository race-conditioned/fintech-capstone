package inbound

// // Result represents the ubiquitous outcome of processing a Command.
// // Any ubiquitous result can be added here.
// type Result interface {
// 	Status() hexa_inbound.ResultStatus
// 	Message() string
// 	Encode(s Sink)
// }
//
// // Sink is an abstraction for writing results to various output mechanisms.
// type Sink interface {
// 	Protocol() string // eg "http", "grpc" etc
// 	Write(status string, v any)
// }

// // ResultStatus represents the status of a transfer operation.
// type ResultStatus string
//
// const (
// 	ResultStatusSuccess     ResultStatus = "success"
// 	ResultStatusRejected    ResultStatus = "rejected"
// 	ResultStatusRateLimited ResultStatus = "rate_limited"
// 	ResultStatusDuplicate   ResultStatus = "duplicate"
// )

// // String returns the string representation of the ResultStatus.
// func (ts ResultStatus) String() string {
// 	return string(ts)
// }
