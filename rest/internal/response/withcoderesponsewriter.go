package response

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

// A WithCodeResponseWriter is a helper to delay sealing a http.ResponseWriter on writing code.
type WithCodeResponseWriter struct {
	Writer http.ResponseWriter
	Code   int
}

// Flush flushes the response writer.
func (w *WithCodeResponseWriter) Flush() {
	if flusher, ok := w.Writer.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Header returns the http header.
func (w *WithCodeResponseWriter) Header() http.Header {
	return w.Writer.Header()
}

// Hijack implements the http.Hijacker interface.
// This expands the Response to fulfill http.Hijacker if the underlying http.ResponseWriter supports it.
func (w *WithCodeResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacked, ok := w.Writer.(http.Hijacker); ok {
		return hijacked.Hijack()
	}

	return nil, nil, errors.New("server doesn't support hijacking")
}

// Write writes bytes into w.
func (w *WithCodeResponseWriter) Write(bytes []byte) (int, error) {
	return w.Writer.Write(bytes)
}

// WriteHeader writes code into w, and not sealing the writer.
func (w *WithCodeResponseWriter) WriteHeader(code int) {
	w.Writer.WriteHeader(code)
	w.Code = code
}
