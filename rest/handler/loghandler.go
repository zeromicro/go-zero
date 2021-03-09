package handler

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/tal-tech/go-zero/core/iox"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/timex"
	"github.com/tal-tech/go-zero/core/utils"
	"github.com/tal-tech/go-zero/rest/httpx"
	"github.com/tal-tech/go-zero/rest/internal"
)

const slowThreshold = time.Millisecond * 500

type loggedResponseWriter struct {
	w    http.ResponseWriter
	r    *http.Request
	code int
}

func (w *loggedResponseWriter) Flush() {
	if flusher, ok := w.w.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *loggedResponseWriter) Header() http.Header {
	return w.w.Header()
}

// Hijack implements the http.Hijacker interface.
// This expands the Response to fulfill http.Hijacker if the underlying http.ResponseWriter supports it.
func (w *loggedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.w.(http.Hijacker).Hijack()
}

func (w *loggedResponseWriter) Write(bytes []byte) (int, error) {
	return w.w.Write(bytes)
}

func (w *loggedResponseWriter) WriteHeader(code int) {
	w.w.WriteHeader(code)
	w.code = code
}

// LogHandler returns a middleware that logs http request and response.
func LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := utils.NewElapsedTimer()
		logs := new(internal.LogCollector)
		lrw := loggedResponseWriter{
			w:    w,
			r:    r,
			code: http.StatusOK,
		}

		var dup io.ReadCloser
		r.Body, dup = iox.DupReadCloser(r.Body)
		next.ServeHTTP(&lrw, r.WithContext(context.WithValue(r.Context(), internal.LogContext, logs)))
		r.Body = dup
		logBrief(r, lrw.code, timer, logs)
	})
}

type detailLoggedResponseWriter struct {
	writer *loggedResponseWriter
	buf    *bytes.Buffer
}

func newDetailLoggedResponseWriter(writer *loggedResponseWriter, buf *bytes.Buffer) *detailLoggedResponseWriter {
	return &detailLoggedResponseWriter{
		writer: writer,
		buf:    buf,
	}
}

func (w *detailLoggedResponseWriter) Flush() {
	w.writer.Flush()
}

func (w *detailLoggedResponseWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *detailLoggedResponseWriter) Write(bs []byte) (int, error) {
	w.buf.Write(bs)
	return w.writer.Write(bs)
}

func (w *detailLoggedResponseWriter) WriteHeader(code int) {
	w.writer.WriteHeader(code)
}

// DetailedLogHandler returns a middleware that logs http request and response in details.
func DetailedLogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := utils.NewElapsedTimer()
		var buf bytes.Buffer
		lrw := newDetailLoggedResponseWriter(&loggedResponseWriter{
			w:    w,
			r:    r,
			code: http.StatusOK,
		}, &buf)

		var dup io.ReadCloser
		r.Body, dup = iox.DupReadCloser(r.Body)
		logs := new(internal.LogCollector)
		next.ServeHTTP(lrw, r.WithContext(context.WithValue(r.Context(), internal.LogContext, logs)))
		r.Body = dup
		logDetails(r, lrw, timer, logs)
	})
}

func dumpRequest(r *http.Request) string {
	reqContent, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err.Error()
	}

	return string(reqContent)
}

func logBrief(r *http.Request, code int, timer *utils.ElapsedTimer, logs *internal.LogCollector) {
	var buf bytes.Buffer
	duration := timer.Duration()
	buf.WriteString(fmt.Sprintf("%d - %s - %s - %s - %s",
		code, r.RequestURI, httpx.GetRemoteAddr(r), r.UserAgent(), timex.ReprOfDuration(duration)))
	if duration > slowThreshold {
		logx.WithContext(r.Context()).Slowf("[HTTP] %d - %s - %s - %s - slowcall(%s)",
			code, r.RequestURI, httpx.GetRemoteAddr(r), r.UserAgent(), timex.ReprOfDuration(duration))
	}

	ok := isOkResponse(code)
	if !ok {
		buf.WriteString(fmt.Sprintf("\n%s", dumpRequest(r)))
	}

	body := logs.Flush()
	if len(body) > 0 {
		buf.WriteString(fmt.Sprintf("\n%s", body))
	}

	if ok {
		logx.WithContext(r.Context()).Info(buf.String())
	} else {
		logx.WithContext(r.Context()).Error(buf.String())
	}
}

func logDetails(r *http.Request, response *detailLoggedResponseWriter, timer *utils.ElapsedTimer,
	logs *internal.LogCollector) {
	var buf bytes.Buffer
	duration := timer.Duration()
	buf.WriteString(fmt.Sprintf("%d - %s - %s\n=> %s\n",
		response.writer.code, r.RemoteAddr, timex.ReprOfDuration(duration), dumpRequest(r)))
	if duration > slowThreshold {
		logx.WithContext(r.Context()).Slowf("[HTTP] %d - %s - slowcall(%s)\n=> %s\n",
			response.writer.code, r.RemoteAddr, timex.ReprOfDuration(duration), dumpRequest(r))
	}

	body := logs.Flush()
	if len(body) > 0 {
		buf.WriteString(fmt.Sprintf("%s\n", body))
	}

	respBuf := response.buf.Bytes()
	if len(respBuf) > 0 {
		buf.WriteString(fmt.Sprintf("<= %s", respBuf))
	}

	logx.WithContext(r.Context()).Info(buf.String())
}

func isOkResponse(code int) bool {
	// not server error
	return code < http.StatusInternalServerError
}
