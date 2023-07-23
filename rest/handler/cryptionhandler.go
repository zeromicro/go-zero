package handler

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"net"
	"net/http"

	"github.com/zeromicro/go-zero/core/codec"
	"github.com/zeromicro/go-zero/core/logx"
)

const maxBytes = 1 << 20 // 1 MiB

var errContentLengthExceeded = errors.New("content length exceeded")

// CryptionHandler returns a middleware to handle cryption.
func CryptionHandler(key []byte) func(http.Handler) http.Handler {
	return LimitCryptionHandler(maxBytes, key)
}

// LimitCryptionHandler returns a middleware to handle cryption.
func LimitCryptionHandler(limitBytes int64, key []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cw := newCryptionResponseWriter(w)
			defer cw.flush(key)

			if r.ContentLength <= 0 {
				next.ServeHTTP(cw, r)
				return
			}

			if err := decryptBody(limitBytes, key, r); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			next.ServeHTTP(cw, r)
		})
	}
}

func decryptBody(limitBytes int64, key []byte, r *http.Request) error {
	if limitBytes > 0 && r.ContentLength > limitBytes {
		return errContentLengthExceeded
	}

	var content []byte
	var err error
	if r.ContentLength > 0 {
		content = make([]byte, r.ContentLength)
		_, err = io.ReadFull(r.Body, content)
	} else {
		content, err = io.ReadAll(io.LimitReader(r.Body, maxBytes))
	}
	if err != nil {
		return err
	}

	content, err = base64.StdEncoding.DecodeString(string(content))
	if err != nil {
		return err
	}

	output, err := codec.EcbDecrypt(key, content)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.Write(output)
	r.Body = io.NopCloser(&buf)

	return nil
}

type cryptionResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func newCryptionResponseWriter(w http.ResponseWriter) *cryptionResponseWriter {
	return &cryptionResponseWriter{
		ResponseWriter: w,
		buf:            new(bytes.Buffer),
	}
}

func (w *cryptionResponseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *cryptionResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Hijack implements the http.Hijacker interface.
// This expands the Response to fulfill http.Hijacker if the underlying http.ResponseWriter supports it.
func (w *cryptionResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacked, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacked.Hijack()
	}

	return nil, nil, errors.New("server doesn't support hijacking")
}

func (w *cryptionResponseWriter) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

func (w *cryptionResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *cryptionResponseWriter) flush(key []byte) {
	if w.buf.Len() == 0 {
		return
	}

	content, err := codec.EcbEncrypt(key, w.buf.Bytes())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body := base64.StdEncoding.EncodeToString(content)
	if n, err := io.WriteString(w.ResponseWriter, body); err != nil {
		logx.Errorf("write response failed, error: %s", err)
	} else if n < len(content) {
		logx.Errorf("actual bytes: %d, written bytes: %d", len(content), n)
	}
}
