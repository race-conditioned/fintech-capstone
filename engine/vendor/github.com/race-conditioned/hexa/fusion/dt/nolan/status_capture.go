package nolan

import "net/http"

// StatusCapture wraps an http.ResponseWriter to capture status code and bytes written.
type statusCapture struct {
	http.ResponseWriter
	status int
	wrote  int
}

// WriteHeader captures the status code.
func (w *statusCapture) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Write captures the number of bytes written.
func (w *statusCapture) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(b)
	w.wrote += n
	return n, err
}
