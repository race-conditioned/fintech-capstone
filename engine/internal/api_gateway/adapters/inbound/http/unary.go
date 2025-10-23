package http_transport

import (
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/platform/apperr"
	"fintech-capstone/m/v2/internal/platform/http_kit/middleware"
	"fintech-capstone/m/v2/internal/platform/http_kit/writer"
	"net/http"
	"strings"
)

type (
	// Decoder decodes an HTTP request into a specific request type.
	Decoder[Req any] func(r *http.Request) (Req, error)
	// Encoder encodes a specific response type into an HTTP response.
	Encoder[Res any] func(w http.ResponseWriter, res Res) // success only
)

// encodeError encodes an error into an HTTP response.
func encodeError(w http.ResponseWriter, err error) {
	e := apperr.As(err)
	switch e.Code {
	case apperr.CodeInvalid:
		writer.JSON(w, http.StatusBadRequest, map[string]string{"error": e.Msg})
	case apperr.CodeRateLimited:
		writer.JSON(w, http.StatusTooManyRequests, map[string]string{"error": e.Msg})
	case apperr.CodeTimeout:
		writer.JSON(w, http.StatusGatewayTimeout, map[string]string{"error": e.Msg})
	case apperr.CodeNotFound:
		writer.JSON(w, http.StatusNotFound, map[string]string{"error": e.Msg})
	case apperr.CodeConflict:
		writer.JSON(w, http.StatusConflict, map[string]string{"error": e.Msg})
	case apperr.CodePayloadTooLarge:
		writer.JSON(w, http.StatusRequestEntityTooLarge, map[string]string{"error": e.Msg})
	default:
		writer.JSON(w, http.StatusInternalServerError, map[string]string{"error": e.Msg})
	}
}

// JSONEncoder returns an Encoder that encodes the response as JSON with the given HTTP status code.
func JSONEncoder[Res any](status int) Encoder[Res] {
	return func(w http.ResponseWriter, res Res) {
		writer.JSON(w, status, res)
	}
}

// EmptyDecoder is a Decoder that returns an empty struct and no error.
func EmptyDecoder(r *http.Request) (struct{}, error) { return struct{}{}, nil }

// Unary creates an HTTP handler function for a unary inbound handler.
func Unary[Req any, Res any](
	h inbound.UnaryHandler[Req, Res],
	dec Decoder[Req],
	enc Encoder[Res],
	metaFrom func(*http.Request) inbound.RequestMeta,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := dec(r)
		if err != nil {
			encodeError(w, err)
			return
		}
		meta := metaFrom(r)
		res, err := h(r.Context(), meta, req)
		if err != nil {
			encodeError(w, err)
			return
		}
		enc(w, res)
	}
}

// DefaultMeta extracts default metadata from the HTTP request.
func DefaultMeta(r *http.Request) inbound.RequestMeta {
	rid, _ := middleware.FromContext(r.Context())
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
