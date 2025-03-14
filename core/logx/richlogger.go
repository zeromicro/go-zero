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
	if shallLog(DebugLevel) {
		l.debug(fmt.Sprint(v...))
	}
}

func (l *richLogger) Debugf(format string, v ...any) {
	if shallLog(DebugLevel) {
		l.debug(fmt.Sprintf(format, v...))
	}
}

func (l *richLogger) Debugfn(fn func() any) {
	if shallLog(DebugLevel) {
		l.debug(fn())
	}
}

func (l *richLogger) Debugv(v any) {
	if shallLog(DebugLevel) {
		l.debug(v)
	}
}

func (l *richLogger) Debugw(msg string, fields ...LogField) {
	if shallLog(DebugLevel) {
		l.debug(msg, fields...)
	}
}

func (l *richLogger) Error(v ...any) {
	if shallLog(ErrorLevel) {
		l.err(fmt.Sprint(v...))
	}
}

func (l *richLogger) Errorf(format string, v ...any) {
	if shallLog(ErrorLevel) {
		l.err(fmt.Sprintf(format, v...))
	}
}

func (l *richLogger) Errorfn(fn func() any) {
	if shallLog(ErrorLevel) {
		l.err(fn())
	}
}

func (l *richLogger) Errorv(v any) {
	if shallLog(ErrorLevel) {
		l.err(v)
	}
}

func (l *richLogger) Errorw(msg string, fields ...LogField) {
	if shallLog(ErrorLevel) {
		l.err(msg, fields...)
	}
}

func (l *richLogger) Info(v ...any) {
	if shallLog(InfoLevel) {
		l.info(fmt.Sprint(v...))
	}
}

func (l *richLogger) Infof(format string, v ...any) {
	if shallLog(InfoLevel) {
		l.info(fmt.Sprintf(format, v...))
	}
}

func (l *richLogger) Infofn(fn func() any) {
	if shallLog(InfoLevel) {
		l.info(fn())
	}
}

func (l *richLogger) Infov(v any) {
	if shallLog(InfoLevel) {
		l.info(v)
	}
}

func (l *richLogger) Infow(msg string, fields ...LogField) {
	if shallLog(InfoLevel) {
		l.info(msg, fields...)
	}
}

func (l *richLogger) Slow(v ...any) {
	if shallLog(ErrorLevel) {
		l.slow(fmt.Sprint(v...))
	}
}

func (l *richLogger) Slowf(format string, v ...any) {
	if shallLog(ErrorLevel) {
		l.slow(fmt.Sprintf(format, v...))
	}
}

func (l *richLogger) Slowfn(fn func() any) {
	if shallLog(ErrorLevel) {
		l.slow(fn())
	}
}

func (l *richLogger) Slowv(v any) {
	if shallLog(ErrorLevel) {
		l.slow(v)
	}
}

func (l *richLogger) Sloww(msg string, fields ...LogField) {
	if shallLog(ErrorLevel) {
		l.slow(msg, fields...)
	}
}

func (l *richLogger) WithCallerSkip(skip int) Logger {
	if skip <= 0 {
		return l
	}

	return &richLogger{
		ctx:        l.ctx,
		callerSkip: skip,
		fields:     l.fields,
	}
}

func (l *richLogger) WithContext(ctx context.Context) Logger {
	return &richLogger{
		ctx:        ctx,
		callerSkip: l.callerSkip,
		fields:     l.fields,
	}
}

func (l *richLogger) WithDuration(duration time.Duration) Logger {
	fields := append(l.fields, Field(durationKey, timex.ReprOfDuration(duration)))

	return &richLogger{
		ctx:        l.ctx,
		callerSkip: l.callerSkip,
		fields:     fields,
	}
}

func (l *richLogger) WithFields(fields ...LogField) Logger {
	if len(fields) == 0 {
		return l
	}

	f := append(l.fields, fields...)

	return &richLogger{
		ctx:        l.ctx,
		callerSkip: l.callerSkip,
		fields:     f,
	}
}

func (l *richLogger) buildFields(fields ...LogField) []LogField {
	fields = append(l.fields, fields...)
	// caller field should always appear together with global fields
	fields = append(fields, Field(callerKey, getCaller(callerDepth+l.callerSkip)))
	fields = mergeGlobalFields(fields)

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
