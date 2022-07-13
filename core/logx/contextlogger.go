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
	return &contextLogger{
		ctx: ctx,
	}
}

type contextLogger struct {
	logEntry
	ctx context.Context
}

func (l *contextLogger) Error(v ...interface{}) {
	l.err(fmt.Sprint(v...))
}

func (l *contextLogger) Errorf(format string, v ...interface{}) {
	l.err(fmt.Sprintf(format, v...))
}

func (l *contextLogger) Errorv(v interface{}) {
	l.err(fmt.Sprint(v))
}

func (l *contextLogger) Errorw(msg string, fields ...LogField) {
	l.err(msg, fields...)
}

func (l *contextLogger) Info(v ...interface{}) {
	l.info(fmt.Sprint(v...))
}

func (l *contextLogger) Infof(format string, v ...interface{}) {
	l.info(fmt.Sprintf(format, v...))
}

func (l *contextLogger) Infov(v interface{}) {
	l.info(v)
}

func (l *contextLogger) Infow(msg string, fields ...LogField) {
	l.info(msg, fields...)
}

func (l *contextLogger) Slow(v ...interface{}) {
	l.slow(fmt.Sprint(v...))
}

func (l *contextLogger) Slowf(format string, v ...interface{}) {
	l.slow(fmt.Sprintf(format, v...))
}

func (l *contextLogger) Slowv(v interface{}) {
	l.slow(v)
}

func (l *contextLogger) Sloww(msg string, fields ...LogField) {
	l.slow(msg, fields...)
}

func (l *contextLogger) WithContext(ctx context.Context) Logger {
	if ctx == nil {
		return l
	}

	l.ctx = ctx
	return l
}

func (l *contextLogger) WithDuration(duration time.Duration) Logger {
	l.Duration = timex.ReprOfDuration(duration)
	return l
}

func (l *contextLogger) buildFields(fields ...LogField) []LogField {
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

	val := l.ctx.Value(fieldsContextKey)
	if val != nil {
		if arr, ok := val.([]LogField); ok {
			fields = append(fields, arr...)
		}
	}

	return fields
}

func (l *contextLogger) err(v interface{}, fields ...LogField) {
	if shallLog(ErrorLevel) {
		getWriter().Error(v, l.buildFields(fields...)...)
	}
}

func (l *contextLogger) info(v interface{}, fields ...LogField) {
	if shallLog(InfoLevel) {
		getWriter().Info(v, l.buildFields(fields...)...)
	}
}

func (l *contextLogger) slow(v interface{}, fields ...LogField) {
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
