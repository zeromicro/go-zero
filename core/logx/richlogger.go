package logx

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/timex"
	"github.com/zeromicro/go-zero/internal/trace"
)

// WithCallerSkip returns a Logger with given caller skip.
func WithCallerSkip(skip int) Logger {
	if skip <= 0 {
		return new(richLogger)
	}

	return &richLogger{
		callerSkip: skip,
	}
}

// WithContext sets ctx to log, for keeping tracing information.
func WithContext(ctx context.Context) Logger {
	return &richLogger{
		ctx: ctx,
	}
}

// WithDuration returns a Logger with given duration.
func WithDuration(d time.Duration) Logger {
	return &richLogger{
		fields: []LogField{Field(durationKey, timex.ReprOfDuration(d))},
	}
}

type richLogger struct {
	ctx        context.Context
	callerSkip int
	fields     []LogField
}

func (l *richLogger) Debug(v ...any) {
	l.debug(fmt.Sprint(v...))
}

func (l *richLogger) Debugf(format string, v ...any) {
	l.debug(fmt.Sprintf(format, v...))
}

func (l *richLogger) Debugv(v any) {
	l.debug(v)
}

func (l *richLogger) Debugw(msg string, fields ...LogField) {
	l.debug(msg, fields...)
}

func (l *richLogger) Error(v ...any) {
	l.err(fmt.Sprint(v...))
}

func (l *richLogger) Errorf(format string, v ...any) {
	l.err(fmt.Sprintf(format, v...))
}

func (l *richLogger) Errorv(v any) {
	l.err(v)
}

func (l *richLogger) Errorw(msg string, fields ...LogField) {
	l.err(msg, fields...)
}

func (l *richLogger) Info(v ...any) {
	l.info(fmt.Sprint(v...))
}

func (l *richLogger) Infof(format string, v ...any) {
	l.info(fmt.Sprintf(format, v...))
}

func (l *richLogger) Infov(v any) {
	l.info(v)
}

func (l *richLogger) Infow(msg string, fields ...LogField) {
	l.info(msg, fields...)
}

func (l *richLogger) Slow(v ...any) {
	l.slow(fmt.Sprint(v...))
}

func (l *richLogger) Slowf(format string, v ...any) {
	l.slow(fmt.Sprintf(format, v...))
}

func (l *richLogger) Slowv(v any) {
	l.slow(v)
}

func (l *richLogger) Sloww(msg string, fields ...LogField) {
	l.slow(msg, fields...)
}

func (l *richLogger) WithCallerSkip(skip int) Logger {
	if skip <= 0 {
		return l
	}

	l.callerSkip = skip
	return l
}

func (l *richLogger) WithContext(ctx context.Context) Logger {
	l.ctx = ctx
	return l
}

func (l *richLogger) WithDuration(duration time.Duration) Logger {
	l.fields = append(l.fields, Field(durationKey, timex.ReprOfDuration(duration)))
	return l
}

func (l *richLogger) WithFields(fields ...LogField) Logger {
	l.fields = append(l.fields, fields...)
	return l
}

func (l *richLogger) buildFields(fields ...LogField) []LogField {
	fields = append(l.fields, fields...)
	fields = append(fields, Field(callerKey, getCaller(callerDepth+l.callerSkip)))

	if l.ctx == nil {
		return fields
	}

	traceID := trace.TraceIDFromContext(l.ctx)
	if len(traceID) > 0 {
		fields = append(fields, Field(traceKey, traceID))
	}

	spanID := trace.SpanIDFromContext(l.ctx)
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

func (l *richLogger) debug(v any, fields ...LogField) {
	if shallLog(DebugLevel) {
		getWriter().Debug(v, l.buildFields(fields...)...)
	}
}

func (l *richLogger) err(v any, fields ...LogField) {
	if shallLog(ErrorLevel) {
		getWriter().Error(v, l.buildFields(fields...)...)
	}
}

func (l *richLogger) info(v any, fields ...LogField) {
	if shallLog(InfoLevel) {
		getWriter().Info(v, l.buildFields(fields...)...)
	}
}

func (l *richLogger) slow(v any, fields ...LogField) {
	if shallLog(ErrorLevel) {
		getWriter().Slow(v, l.buildFields(fields...)...)
	}
}
