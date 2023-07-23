package internal

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// logContextKey is a context key.
var logContextKey = contextKey("request_logs")

type (
	// LogCollector is used to collect logs.
	LogCollector struct {
		Messages []string
		lock     sync.Mutex
	}

	contextKey string
)

// WithLogCollector returns a new context with LogCollector.
func WithLogCollector(ctx context.Context, lc *LogCollector) context.Context {
	return context.WithValue(ctx, logContextKey, lc)
}

// LogCollectorFromContext returns LogCollector from ctx.
func LogCollectorFromContext(ctx context.Context) *LogCollector {
	val := ctx.Value(logContextKey)
	if val == nil {
		return nil
	}

	return val.(*LogCollector)
}

// Append appends msg into log context.
func (lc *LogCollector) Append(msg string) {
	lc.lock.Lock()
	lc.Messages = append(lc.Messages, msg)
	lc.lock.Unlock()
}

// Flush flushes collected logs.
func (lc *LogCollector) Flush() string {
	var buffer bytes.Buffer

	start := true
	for _, message := range lc.takeAll() {
		if start {
			start = false
		} else {
			buffer.WriteByte('\n')
		}
		buffer.WriteString(message)
	}

	return buffer.String()
}

func (lc *LogCollector) takeAll() []string {
	lc.lock.Lock()
	messages := lc.Messages
	lc.Messages = nil
	lc.lock.Unlock()

	return messages
}

// Error logs the given v along with r in error log.
func Error(r *http.Request, v ...any) {
	logx.WithContext(r.Context()).Error(format(r, v...))
}

// Errorf logs the given v with format along with r in error log.
func Errorf(r *http.Request, format string, v ...any) {
	logx.WithContext(r.Context()).Error(formatf(r, format, v...))
}

// Info logs the given v along with r in access log.
func Info(r *http.Request, v ...any) {
	appendLog(r, format(r, v...))
}

// Infof logs the given v with format along with r in access log.
func Infof(r *http.Request, format string, v ...any) {
	appendLog(r, formatf(r, format, v...))
}

func appendLog(r *http.Request, message string) {
	logs := LogCollectorFromContext(r.Context())
	if logs != nil {
		logs.Append(message)
	}
}

func format(r *http.Request, v ...any) string {
	return formatWithReq(r, fmt.Sprint(v...))
}

func formatf(r *http.Request, format string, v ...any) string {
	return formatWithReq(r, fmt.Sprintf(format, v...))
}

func formatWithReq(r *http.Request, v string) string {
	return fmt.Sprintf("(%s - %s) %s", r.RequestURI, httpx.GetRemoteAddr(r), v)
}
