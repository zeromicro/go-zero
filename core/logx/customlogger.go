package logx

import (
	"fmt"
	"io"
	"time"

	"github.com/tal-tech/go-zero/core/timex"
)

const customCallerDepth = 3

type customLog logEntry

func WithDuration(d time.Duration) Logger {
	return customLog{
		Duration: timex.ReprOfDuration(d),
	}
}

func (l customLog) Error(v ...interface{}) {
	if shouldLog(ErrorLevel) {
		l.write(errorLog, levelError, formatWithCaller(fmt.Sprint(v...), customCallerDepth))
	}
}

func (l customLog) Errorf(format string, v ...interface{}) {
	if shouldLog(ErrorLevel) {
		l.write(errorLog, levelError, formatWithCaller(fmt.Sprintf(format, v...), customCallerDepth))
	}
}

func (l customLog) Info(v ...interface{}) {
	if shouldLog(InfoLevel) {
		l.write(infoLog, levelInfo, fmt.Sprint(v...))
	}
}

func (l customLog) Infof(format string, v ...interface{}) {
	if shouldLog(InfoLevel) {
		l.write(infoLog, levelInfo, fmt.Sprintf(format, v...))
	}
}

func (l customLog) Slow(v ...interface{}) {
	if shouldLog(ErrorLevel) {
		l.write(slowLog, levelSlow, fmt.Sprint(v...))
	}
}

func (l customLog) Slowf(format string, v ...interface{}) {
	if shouldLog(ErrorLevel) {
		l.write(slowLog, levelSlow, fmt.Sprintf(format, v...))
	}
}

func (l customLog) write(writer io.Writer, level, content string) {
	l.Timestamp = getTimestamp()
	l.Level = level
	l.Content = content
	outputJson(writer, logEntry(l))
}
