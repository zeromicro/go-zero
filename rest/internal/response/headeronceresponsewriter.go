package response

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

// HeaderOnceResponseWriter is a http.ResponseWriter implementation
// that only the first WriterHeader takes effect.
type HeaderOnceResponseWriter struct {
	w           http.ResponseWriter
	wroteHeader bool
}

// NewHeaderOnceResponseWriter returns a HeaderOnceResponseWriter.
func NewHeaderOnceResponseWriter(w http.ResponseWriter) http.ResponseWriter {
	return &HeaderOnceResponseWriter{w: w}
}

// Flush flushes the response writer.
func (w *HeaderOnceResponseWriter) Flush() {
	if flusher, ok := w.w.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Header returns the http header.
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

// Write writes bytes into w.
func (w *HeaderOnceResponseWriter) Write(bytes []byte) (int, error) {
	return w.w.Write(bytes)
}

// WriteHeader writes code into w, and not sealing the writer.
func (w *HeaderOnceResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}

	w.w.WriteHeader(code)
	w.wroteHeader = true
}
