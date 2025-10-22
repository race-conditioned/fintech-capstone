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
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(buf); err != nil {
		log.Printf("writer.JSON: write failed: %v", err)
	}
}

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]string{"error": msg})
}
