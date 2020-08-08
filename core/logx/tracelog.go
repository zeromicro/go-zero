package logx

import (
	"context"
	"fmt"
	"io"

	"github.com/tal-tech/go-zero/core/trace/tracespec"
)

type tracingEntry struct {
	logEntry
	Trace string          `json:"trace,omitempty"`
	Span  string          `json:"span,omitempty"`
	ctx   context.Context `json:"-"`
}

func (l tracingEntry) Error(v ...interface{}) {
	if shouldLog(ErrorLevel) {
		l.write(errorLog, levelError, formatWithCaller(fmt.Sprint(v...), customCallerDepth))
	}
}

func (l tracingEntry) Errorf(format string, v ...interface{}) {
	if shouldLog(ErrorLevel) {
		l.write(errorLog, levelError, formatWithCaller(fmt.Sprintf(format, v...), customCallerDepth))
	}
}

func (l tracingEntry) Info(v ...interface{}) {
	if shouldLog(InfoLevel) {
		l.write(infoLog, levelInfo, fmt.Sprint(v...))
	}
}

func (l tracingEntry) Infof(format string, v ...interface{}) {
	if shouldLog(InfoLevel) {
		l.write(infoLog, levelInfo, fmt.Sprintf(format, v...))
	}
}

func (l tracingEntry) Slow(v ...interface{}) {
	if shouldLog(ErrorLevel) {
		l.write(slowLog, levelSlow, fmt.Sprint(v...))
	}
}

func (l tracingEntry) Slowf(format string, v ...interface{}) {
	if shouldLog(ErrorLevel) {
		l.write(slowLog, levelSlow, fmt.Sprintf(format, v...))
	}
}

func (l tracingEntry) write(writer io.Writer, level, content string) {
	l.Timestamp = getTimestamp()
	l.Level = level
	l.Content = content
	l.Trace = traceIdFromContext(l.ctx)
	l.Span = spanIdFromContext(l.ctx)
	outputJson(writer, l)
}

func WithContext(ctx context.Context) Logger {
	return tracingEntry{
		ctx: ctx,
	}
}

func spanIdFromContext(ctx context.Context) string {
	t, ok := ctx.Value(tracespec.TracingKey).(tracespec.Trace)
	if !ok {
		return ""
	}

	return t.SpanId()
}

func traceIdFromContext(ctx context.Context) string {
	t, ok := ctx.Value(tracespec.TracingKey).(tracespec.Trace)
	if !ok {
		return ""
	}

	return t.TraceId()
}
