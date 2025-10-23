package apperr

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

type Error struct {
	Code Code
	Msg  string
	Err  error // optional wrap
}

func (e *Error) Error() string { return e.Msg }

func Wrap(code Code, msg string, err error) *Error { return &Error{Code: code, Msg: msg, Err: err} }

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
