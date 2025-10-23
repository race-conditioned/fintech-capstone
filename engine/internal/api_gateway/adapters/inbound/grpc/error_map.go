package grpc_transport

import (
	"errors"
	"fintech-capstone/m/v2/internal/platform/apperr"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// toGRPCError converts an application error to a gRPC status error.
func toGRPCError(err error) error {
	if err == nil {
		return nil
	}

	e := apperr.As(err)
	var code codes.Code

	switch e.Code {
	case apperr.CodeInvalid:
		code = codes.InvalidArgument
	case apperr.CodeRateLimited:
		code = codes.ResourceExhausted
	case apperr.CodeTimeout:
		code = codes.DeadlineExceeded
	case apperr.CodeNotFound:
		code = codes.NotFound
	case apperr.CodeConflict:
		code = codes.AlreadyExists // clearer than Aborted for idempotency
	case apperr.CodePayloadTooLarge:
		code = codes.ResourceExhausted
	default:
		code = codes.Internal
	}

	// Preserve root cause if present
	if e.Err != nil && !errors.Is(err, e.Err) {
		return status.Errorf(code, "%s: %v", e.Msg, e.Err)
	}
	return status.Error(code, e.Msg)
}
