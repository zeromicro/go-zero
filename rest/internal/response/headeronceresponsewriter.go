package response

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

type HeaderOnceResponseWriter struct {
	w           http.ResponseWriter
	wroteHeader bool
}

func NewHeaderOnceResponseWriter(w http.ResponseWriter) http.ResponseWriter {
	return &HeaderOnceResponseWriter{w: w}
}

func (w *HeaderOnceResponseWriter) Flush() {
	if flusher, ok := w.w.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *HeaderOnceResponseWriter) Header() http.Header {
	return w.w.Header()
}

// Hijack implements the http.Hijacker interface.
// This expands the Response to fulfill http.Hijacker if the underlying http.ResponseWriter supports it.
func (w *HeaderOnceResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacked, ok := w.w.(http.Hijacker); ok {
		return hijacked.Hijack()
	}

	return nil, nil, errors.New("server doesn't support hijacking")
}

func (w *HeaderOnceResponseWriter) Write(bytes []byte) (int, error) {
	return w.w.Write(bytes)
}

func (w *HeaderOnceResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}

	w.w.WriteHeader(code)
	w.wroteHeader = true
}
