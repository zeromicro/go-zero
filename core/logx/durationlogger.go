package logx

import (
	"fmt"
	"io"
	"time"

	"github.com/zeromicro/go-zero/core/timex"
	"golang.org/x/net/context"
)

type durationLogger logEntry

// WithDuration returns a Logger which logs the given duration.
func WithDuration(d time.Duration) Logger {
	return &durationLogger{
		Duration: timex.ReprOfDuration(d),
	}
}

func (l *durationLogger) Error(v ...interface{}) {
	if shallLog(ErrorLevel) {
		l.write(errorLog, levelError, fmt.Sprint(v...))
	}
}

func (l *durationLogger) Errorf(format string, v ...interface{}) {
	if shallLog(ErrorLevel) {
		l.write(errorLog, levelError, fmt.Sprintf(format, v...))
	}
}

func (l *durationLogger) Errorv(v interface{}) {
	if shallLog(ErrorLevel) {
		l.write(errorLog, levelError, v)
	}
}

func (l *durationLogger) Errorw(msg string, fields ...LogField) {
	if shallLog(ErrorLevel) {
		l.write(errorLog, levelError, msg, fields...)
	}
}

func (l *durationLogger) Info(v ...interface{}) {
	if shallLog(InfoLevel) {
		l.write(infoLog, levelInfo, fmt.Sprint(v...))
	}
}

func (l *durationLogger) Infof(format string, v ...interface{}) {
	if shallLog(InfoLevel) {
		l.write(infoLog, levelInfo, fmt.Sprintf(format, v...))
	}
}

func (l *durationLogger) Infov(v interface{}) {
	if shallLog(InfoLevel) {
		l.write(infoLog, levelInfo, v)
	}
}

func (l *durationLogger) Infow(msg string, fields ...LogField) {
	if shallLog(InfoLevel) {
		l.write(infoLog, levelInfo, msg, fields...)
	}
}

func (l *durationLogger) Slow(v ...interface{}) {
	if shallLog(ErrorLevel) {
		l.write(slowLog, levelSlow, fmt.Sprint(v...))
	}
}

func (l *durationLogger) Slowf(format string, v ...interface{}) {
	if shallLog(ErrorLevel) {
		l.write(slowLog, levelSlow, fmt.Sprintf(format, v...))
	}
}

func (l *durationLogger) Slowv(v interface{}) {
	if shallLog(ErrorLevel) {
		l.write(slowLog, levelSlow, v)
	}
}

func (l *durationLogger) Sloww(msg string, fields ...LogField) {
	if shallLog(ErrorLevel) {
		l.write(slowLog, levelSlow, msg, fields...)
	}
}

func (l *durationLogger) WithContext(ctx context.Context) Logger {
	return &traceLogger{
		ctx: ctx,
		logEntry: logEntry{
			Duration: l.Duration,
		},
	}
}

func (l *durationLogger) WithDuration(duration time.Duration) Logger {
	l.Duration = timex.ReprOfDuration(duration)
	return l
}

func (l *durationLogger) write(writer io.Writer, level string, val interface{}, fields ...LogField) {
	fields = append(fields, Field(durationKey, l.Duration))
	output(writer, level, val, fields...)
}
