package nolan

import (
	"encoding/json"
	"log"
	"net/http"
)

type JSONWriter interface {
	Encode(w http.ResponseWriter, status int, v any)
}

type Writer struct{}

func NewWriter() *Writer {
	return &Writer{}
}

// Encode writes a JSON response with the given status code.
// If encoding fails, it logs and falls back to an internal error payload.
func (writer *Writer) Encode(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")

	buf, err := json.Marshal(v)
	if err != nil {
		log.Printf("writer.JSON: encode failed: %v", err)
		// Fall back to a generic error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"internal server error"}`))
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(buf); err != nil {
		log.Printf("writer.JSON: write failed: %v", err)
	}
}

// Error writes a JSON error response with the given status code and message.
func (writer *Writer) Error(w http.ResponseWriter, status int, msg string) {
	writer.Encode(w, status, map[string]string{"error": msg})
}

// Wrap wraps an http.ResponseWriter to capture status code and bytes written.
func (writer *Writer) Wrap(w http.ResponseWriter) *statusCapture {
	return &statusCapture{ResponseWriter: w}
}

// Status retrieves the captured status code and bytes written.
func (writer *Writer) Status(w http.ResponseWriter) (code, bytes int, ok bool) {
	if sw, ok2 := w.(*statusCapture); ok2 {
		return sw.status, sw.wrote, true
	}
	return 0, 0, false
}
