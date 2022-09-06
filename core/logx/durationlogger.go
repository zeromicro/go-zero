package logx

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/timex"
)

// WithDuration returns a Logger which logs the given duration.
func WithDuration(d time.Duration) Logger {
	return &durationLogger{
		logEntry:          logEntry{Duration: timex.ReprOfDuration(d)},
		defaultCallerSkip: 4,
	}
}

type durationLogger struct {
	logEntry
	defaultCallerSkip int
}

func (l *durationLogger) Error(v ...interface{}) {
	l.err(fmt.Sprint(v...))
}

func (l *durationLogger) Errorf(format string, v ...interface{}) {
	l.err(fmt.Sprintf(format, v...))
}

func (l *durationLogger) Errorv(v interface{}) {
	l.err(v)
}

func (l *durationLogger) Errorw(msg string, fields ...LogField) {
	l.err(msg, fields...)
}

func (l *durationLogger) Info(v ...interface{}) {
	l.info(fmt.Sprint(v...))
}

func (l *durationLogger) Infof(format string, v ...interface{}) {
	l.info(fmt.Sprintf(format, v...))
}

func (l *durationLogger) Infov(v interface{}) {
	l.info(v)
}

func (l *durationLogger) Infow(msg string, fields ...LogField) {
	l.info(msg, fields...)
}

func (l *durationLogger) Slow(v ...interface{}) {
	l.slow(fmt.Sprint(v...))
}

func (l *durationLogger) Slowf(format string, v ...interface{}) {
	l.slow(fmt.Sprintf(format, v...))
}

func (l *durationLogger) Slowv(v interface{}) {
	l.slow(v)
}

func (l *durationLogger) Sloww(msg string, fields ...LogField) {
	l.slow(msg, fields...)
}

func (l *durationLogger) WithContext(ctx context.Context) Logger {
	return &contextLogger{
		ctx: ctx,
		logEntry: logEntry{
			Duration: l.Duration,
		},
	}
}

func (l *durationLogger) WithCallerDepth(callerDepth int) Logger {
	l.CallerDepth = callerDepth
	return l
}

func (l *durationLogger) WithDuration(duration time.Duration) Logger {
	l.Duration = timex.ReprOfDuration(duration)
	return l
}

func (l *durationLogger) err(v interface{}, fields ...LogField) {
	if shallLog(ErrorLevel) {
		getWriter().Error(v, l.buildFields(fields...)...)
	}
}

func (l *durationLogger) info(v interface{}, fields ...LogField) {
	if shallLog(InfoLevel) {
		getWriter().Info(v, l.buildFields(fields...)...)
	}
}

func (l *durationLogger) slow(v interface{}, fields ...LogField) {
	if shallLog(ErrorLevel) {
		getWriter().Slow(v, l.buildFields(fields...)...)
	}
}

func (l *durationLogger) buildFields(fields ...LogField) []LogField {
	return append(fields, Field(durationKey, l.Duration),
		Field(callerKey, getCaller(l.defaultCallerSkip+l.CallerDepth)))
}
