// Package drain provides an HTTP middleware that asks clients to close their
// connections while the server is draining before a shutdown.
package drain

import (
	"fmt"
	"io"
	"net/http"
)

type connectionCloseWriter struct {
	http.ResponseWriter

	isDraining    func() bool
	headerWritten bool
}

func (w *connectionCloseWriter) WriteHeader(code int) {
	w.injectHeader()
	w.ResponseWriter.WriteHeader(code)
}

func (w *connectionCloseWriter) Write(b []byte) (int, error) {
	// Write implicitly sends a 200 WriteHeader if not yet called,
	// so we inject before that happens.
	w.injectHeader()
	n, err := w.ResponseWriter.Write(b)

	return n, wrapWriteError(err)
}

func (w *connectionCloseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *connectionCloseWriter) ReadFrom(src io.Reader) (int64, error) {
	w.injectHeader()
	if rf, ok := w.ResponseWriter.(io.ReaderFrom); ok {
		n, err := rf.ReadFrom(src)

		return n, wrapWriteError(err)
	}
	// Fallback: copy manually via Write (which won't re-inject thanks to the guard).
	n, err := io.Copy(w.ResponseWriter, src)

	return n, wrapWriteError(err)
}

func (w *connectionCloseWriter) injectHeader() {
	if !w.headerWritten {
		w.headerWritten = true
		if w.isDraining() {
			w.ResponseWriter.Header().Set("Connection", "close")
		}
	}
}

// wrapWriteError adds context to an error from writing the response body.
func wrapWriteError(err error) error {
	if err != nil {
		return fmt.Errorf("writing response body: %w", err)
	}

	return nil
}

// WithConnectionDraining returns a middleware that adds a "Connection: close"
// header to every response while isDraining reports true, so that clients
// stop reusing their connections to this server.
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
