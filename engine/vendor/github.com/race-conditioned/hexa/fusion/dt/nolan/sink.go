package nolan

import (
	"net/http"
)

type Sink struct {
	w    http.ResponseWriter
	json *Writer
}

func NewSink(w http.ResponseWriter) *Sink {
	return &Sink{w: w, json: NewWriter()}
}

func (s *Sink) Protocol() string { return "http" }

func (s *Sink) Write(status string, v any) {
	s.json.Encode(s.w, toHTTP(status), v)
}

func toHTTP(status string) int {
	switch status {
	case "success":
		return http.StatusOK
	case "rejected":
		return http.StatusBadRequest
	// ...
	default:
		return http.StatusInternalServerError
	}
}
