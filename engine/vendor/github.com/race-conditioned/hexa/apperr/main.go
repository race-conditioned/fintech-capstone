package apperr

// Code defines the type for error codes for consistent mapping to transport layers.
type Code int

const (
	CodeOK Code = iota
	CodeInvalid
	CodeRateLimited
	CodeTimeout
	CodeNotFound
	CodePayloadTooLarge
	CodeConflict
	CodeInternal
)

// Error represents a standard application error with a code and message.
type Error struct {
	Code Code
	Msg  string
	Err  error // optional wrap
}

// Error implements the error interface.
func (e *Error) Error() string { return e.Msg }

// Wrap wraps an existing error with a Code and message.
func Wrap(code Code, msg string, err error) *Error { return &Error{Code: code, Msg: msg, Err: err} }

// As converts a generic error to an *Error.
func As(err error) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return &Error{Code: CodeInternal, Msg: err.Error(), Err: err}
}

// Small constructors so adapters/policy can be expressive.
func Invalid(msg string) *Error         { return &Error{Code: CodeInvalid, Msg: msg} }
func RateLimited(msg string) *Error     { return &Error{Code: CodeRateLimited, Msg: msg} }
func Timeout(msg string) *Error         { return &Error{Code: CodeTimeout, Msg: msg} }
func Conflict(msg string) *Error        { return &Error{Code: CodeConflict, Msg: msg} }
func Internal(msg string) *Error        { return &Error{Code: CodeInternal, Msg: msg} }
func PayloadTooLarge(msg string) *Error { return &Error{Code: CodePayloadTooLarge, Msg: msg} }
