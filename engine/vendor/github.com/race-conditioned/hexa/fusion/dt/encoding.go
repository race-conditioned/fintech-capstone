package dt

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/race-conditioned/hexa/apperr"
	"github.com/race-conditioned/hexa/fusion/dt/nolan"
	"github.com/race-conditioned/hexa/horizon/ports/inbound"
)

type (
	// Decoder decodes an HTTP request into a specific request type.
	Decoder[Req any] func(r *http.Request, command *Req) error
	// Encoder encodes a specific response type into an HTTP response.
	Encoder[Res any] func(w http.ResponseWriter, res Res) // success only
)

// Unary creates an HTTP handler when the request type is only known at runtime.
func Unary[c inbound.Ctx](
	handler inbound.UnaryHandler[c, inbound.Command, inbound.Result],
	newPayload func() any,
	metaFrom func(*http.Request) inbound.RequestMeta,
	ctx c,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("trace: Unary")
		// 1) allocate the right concrete command
		payload := newPayload()

		fmt.Println("trace: allocated payload", payload)
		fmt.Printf("trace: allocated payload of type %T\n", payload)

		fmt.Println("body:", r.Body)
		// 2) decode JSON into it
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(payload); err != nil {
			fmt.Println("error here", err)
			encodeError(w, err)
			return
		}

		fmt.Printf("trace: decoded payload: %+v\n", payload)

		// 3) figure out what to pass to the gateway
		var cmd inbound.Command

		switch v := any(payload).(type) {
		case inbound.CommandDTO:
			cmd = v.ToCommand()
		case inbound.Command:
			cmd = v
		default:
			encodeError(w, apperr.Invalid("payload is neither Command nor CommandDTO"))
			return
		}
		fmt.Printf("trace: constructed command of type %T\n", cmd)

		// 4) build meta
		meta := metaFrom(r)
		fmt.Printf("trace: constructed meta: %+v\n", meta)

		// 5) call the gateway handler
		res, err := handler(ctx, meta, cmd)
		if err != nil {
			encodeError(w, err)
			return
		}

		fmt.Printf("trace: obtained result of type %T\n", res)

		// 6) write via sink
		sink := nolan.NewSink(w)
		res.Encode(sink)
		fmt.Println("trace: response encoded successfully")
	}
}

// JSONDecoder decodes a TransferCommand from a JSON HTTP request.
func JSONDecoder[Req any]() Decoder[Req] {
	return func(r *http.Request, command *Req) error {
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&command); err != nil {
			var maxErr *http.MaxBytesError
			if errors.As(err, &maxErr) {
				return apperr.PayloadTooLarge("request body too large")
			}
			return apperr.Invalid("invalid JSON payload")
		}
		return nil
	}
}

// DefaultMeta extracts default metadata from the HTTP request.
func DefaultMeta(r *http.Request) inbound.RequestMeta {
	rid, _ := FromContext(r.Context())
	return inbound.RequestMeta{
		ClientID:  r.Header.Get("X-Client-ID"),
		RequestID: firstNonEmpty(rid, r.Header.Get("X-Request-ID")),
		TraceID:   r.Header.Get("X-Trace-ID"),
		RemoteIP:  r.RemoteAddr,
		Protocol:  "http",
		Target:    r.Method + " " + r.URL.Path,
	}
}

// firstNonEmpty returns the first non-empty string from the provided values.
func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

// encodeError encodes an error into an HTTP response.
func encodeError(w http.ResponseWriter, err error) {
	e := apperr.As(err)
	writer := nolan.NewWriter()
	switch e.Code {
	case apperr.CodeInvalid:
		writer.Encode(w, http.StatusBadRequest, map[string]string{"error": e.Msg})
	case apperr.CodeRateLimited:
		writer.Encode(w, http.StatusTooManyRequests, map[string]string{"error": e.Msg})
	case apperr.CodeTimeout:
		writer.Encode(w, http.StatusGatewayTimeout, map[string]string{"error": e.Msg})
	case apperr.CodeNotFound:
		writer.Encode(w, http.StatusNotFound, map[string]string{"error": e.Msg})
	case apperr.CodeConflict:
		writer.Encode(w, http.StatusConflict, map[string]string{"error": e.Msg})
	case apperr.CodePayloadTooLarge:
		writer.Encode(w, http.StatusRequestEntityTooLarge, map[string]string{"error": e.Msg})
	default:
		writer.Encode(w, http.StatusInternalServerError, map[string]string{"error": e.Msg})
	}
}
