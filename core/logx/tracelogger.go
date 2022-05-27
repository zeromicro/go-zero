package logx

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/timex"
	"go.opentelemetry.io/otel/trace"
)

// WithContext sets ctx to log, for keeping tracing information.
func WithContext(ctx context.Context) Logger {
	return &traceLogger{
		ctx: ctx,
	}
}

type traceLogger struct {
	logEntry
	ctx context.Context
}

func (l *traceLogger) Error(v ...interface{}) {
	l.err(fmt.Sprint(v...))
}

func (l *traceLogger) Errorf(format string, v ...interface{}) {
	l.err(fmt.Sprintf(format, v...))
}

func (l *traceLogger) Errorv(v interface{}) {
	l.err(fmt.Sprint(v))
}

func (l *traceLogger) Errorw(msg string, fields ...LogField) {
	l.err(msg, fields...)
}

func (l *traceLogger) Info(v ...interface{}) {
	l.info(fmt.Sprint(v...))
}

func (l *traceLogger) Infof(format string, v ...interface{}) {
	l.info(fmt.Sprintf(format, v...))
}

func (l *traceLogger) Infov(v interface{}) {
	l.info(v)
}

func (l *traceLogger) Infow(msg string, fields ...LogField) {
	l.info(msg, fields...)
}

func (l *traceLogger) Slow(v ...interface{}) {
	l.slow(fmt.Sprint(v...))
}

func (l *traceLogger) Slowf(format string, v ...interface{}) {
	l.slow(fmt.Sprintf(format, v...))
}

func (l *traceLogger) Slowv(v interface{}) {
	l.slow(v)
}

func (l *traceLogger) Sloww(msg string, fields ...LogField) {
	l.slow(msg, fields...)
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

func (l *traceLogger) buildFields(fields ...LogField) []LogField {
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

	return fields
}

func (l *traceLogger) err(v interface{}, fields ...LogField) {
	if shallLog(ErrorLevel) {
		getWriter().Error(v, l.buildFields(fields...)...)
	}
}

func (l *traceLogger) info(v interface{}, fields ...LogField) {
	if shallLog(InfoLevel) {
		getWriter().Info(v, l.buildFields(fields...)...)
	}
}

func (l *traceLogger) slow(v interface{}, fields ...LogField) {
	if shallLog(ErrorLevel) {
		getWriter().Slow(v, l.buildFields(fields...)...)
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
