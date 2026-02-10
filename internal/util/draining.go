package util

import (
	"io"
	"net/http"
)

type connectionCloseWriter struct {
	http.ResponseWriter
	headerWritten bool
	isDraining    func() bool
}

func (w *connectionCloseWriter) injectHeader() {
	if !w.headerWritten {
		w.headerWritten = true
		if w.isDraining() {
			w.ResponseWriter.Header().Set("Connection", "close")
		}
	}
}

func (w *connectionCloseWriter) WriteHeader(code int) {
	w.injectHeader()
	w.ResponseWriter.WriteHeader(code)
}

func (w *connectionCloseWriter) Write(b []byte) (int, error) {
	// Write implicitly sends a 200 WriteHeader if not yet called,
	// so we inject before that happens.
	w.injectHeader()
	return w.ResponseWriter.Write(b)
}

func (w *connectionCloseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *connectionCloseWriter) ReadFrom(src io.Reader) (int64, error) {
	w.injectHeader()
	if rf, ok := w.ResponseWriter.(io.ReaderFrom); ok {
		return rf.ReadFrom(src)
	}
	// Fallback: copy manually via Write (which won't re-inject thanks to the guard).
	return io.Copy(w.ResponseWriter, src)
}

func WithConnectionDraining(isDraining func() bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cw := &connectionCloseWriter{
				ResponseWriter: w,
				isDraining:     isDraining,
			}
			next.ServeHTTP(cw, r)
		})
	}
}
