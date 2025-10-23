package writer

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON writes a JSON response with the given status code.
// If encoding fails, it logs and falls back to an internal error payload.
func JSON(w http.ResponseWriter, status int, v any) {
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
func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]string{"error": msg})
}
