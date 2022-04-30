package logx

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/zeromicro/go-zero/core/timex"
	"go.opentelemetry.io/otel/trace"
)

type traceLogger struct {
	logEntry
	ctx context.Context
}

func (l *traceLogger) Error(v ...interface{}) {
	if shallLog(ErrorLevel) {
		l.write(errorLog, levelError, fmt.Sprint(v...))
	}
}

func (l *traceLogger) Errorf(format string, v ...interface{}) {
	if shallLog(ErrorLevel) {
		l.write(errorLog, levelError, fmt.Sprintf(format, v...))
	}
}

func (l *traceLogger) Errorv(v interface{}) {
	if shallLog(ErrorLevel) {
		l.write(errorLog, levelError, v)
	}
}

func (l *traceLogger) Errorw(msg string, fields ...LogField) {
	if shallLog(ErrorLevel) {
		l.write(errorLog, levelError, msg, fields...)
	}
}

func (l *traceLogger) Info(v ...interface{}) {
	if shallLog(InfoLevel) {
		l.write(infoLog, levelInfo, fmt.Sprint(v...))
	}
}

func (l *traceLogger) Infof(format string, v ...interface{}) {
	if shallLog(InfoLevel) {
		l.write(infoLog, levelInfo, fmt.Sprintf(format, v...))
	}
}

func (l *traceLogger) Infov(v interface{}) {
	if shallLog(InfoLevel) {
		l.write(infoLog, levelInfo, v)
	}
}

func (l *traceLogger) Infow(msg string, fields ...LogField) {
	if shallLog(InfoLevel) {
		l.write(infoLog, levelInfo, msg, fields...)
	}
}

func (l *traceLogger) Slow(v ...interface{}) {
	if shallLog(ErrorLevel) {
		l.write(slowLog, levelSlow, fmt.Sprint(v...))
	}
}

func (l *traceLogger) Slowf(format string, v ...interface{}) {
	if shallLog(ErrorLevel) {
		l.write(slowLog, levelSlow, fmt.Sprintf(format, v...))
	}
}

func (l *traceLogger) Slowv(v interface{}) {
	if shallLog(ErrorLevel) {
		l.write(slowLog, levelSlow, v)
	}
}

func (l *traceLogger) Sloww(msg string, fields ...LogField) {
	if shallLog(ErrorLevel) {
		l.write(slowLog, levelSlow, msg, fields...)
	}
}

func (l *traceLogger) WithContext(ctx context.Context) Logger {
	if ctx == nil {
		return l
	}

	l.ctx = ctx
	return l
}

func (l *traceLogger) WithDuration(duration time.Duration) Logger {
	l.Duration = timex.ReprOfDuration(duration)
	return l
}

func (l *traceLogger) write(writer io.Writer, level string, val interface{}, fields ...LogField) {
	if len(l.Duration) > 0 {
		fields = append(fields, Field(durationKey, l.Duration))
	}
	traceID := traceIdFromContext(l.ctx)
	if len(traceID) > 0 {
		fields = append(fields, Field(traceKey, traceID))
	}
	spanID := spanIdFromContext(l.ctx)
	if len(spanID) > 0 {
		fields = append(fields, Field(spanKey, spanID))
	}
	outputAny(writer, level, val, fields...)
}

// WithContext sets ctx to log, for keeping tracing information.
func WithContext(ctx context.Context) Logger {
	return &traceLogger{
		ctx: ctx,
	}
}

func spanIdFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		return spanCtx.SpanID().String()
	}

	return ""
}

func traceIdFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}

	return ""
}
