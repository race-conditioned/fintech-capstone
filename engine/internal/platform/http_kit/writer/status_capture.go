package writer

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

// Wrap wraps an http.ResponseWriter to capture status code and bytes written.
func Wrap(w http.ResponseWriter) *statusCapture {
	return &statusCapture{ResponseWriter: w}
}

// Status retrieves the captured status code and bytes written.
func Status(w http.ResponseWriter) (code, bytes int, ok bool) {
	if sw, ok2 := w.(*statusCapture); ok2 {
		return sw.status, sw.wrote, true
	}
	return 0, 0, false
}
