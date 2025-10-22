package writer

import "net/http"

type statusCapture struct {
	http.ResponseWriter
	status int
	wrote  int
}

func (w *statusCapture) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusCapture) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(b)
	w.wrote += n
	return n, err
}

func Wrap(w http.ResponseWriter) *statusCapture {
	return &statusCapture{ResponseWriter: w}
}

func Status(w http.ResponseWriter) (code, bytes int, ok bool) {
	if sw, ok2 := w.(*statusCapture); ok2 {
		return sw.status, sw.wrote, true
	}
	return 0, 0, false
}
