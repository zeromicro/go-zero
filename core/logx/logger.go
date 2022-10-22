package logx

import (
	"context"
	"time"
)

// A Logger represents a logger.
type Logger interface {
	// Debug logs a message at info level.
	Debug(...interface{})
	// Debugf logs a message at info level.
	Debugf(string, ...interface{})
	// Debugv logs a message at info level.
	Debugv(interface{})
	// Debugw logs a message at info level.
	Debugw(string, ...LogField)
	// Error logs a message at error level.
	Error(...interface{})
	// Errorf logs a message at error level.
	Errorf(string, ...interface{})
	// Errorv logs a message at error level.
	Errorv(interface{})
	// Errorw logs a message at error level.
	Errorw(string, ...LogField)
	// Info logs a message at info level.
	Info(...interface{})
	// Infof logs a message at info level.
	Infof(string, ...interface{})
	// Infov logs a message at info level.
	Infov(interface{})
	// Infow logs a message at info level.
	Infow(string, ...LogField)
	// Slow logs a message at slow level.
	Slow(...interface{})
	// Slowf logs a message at slow level.
	Slowf(string, ...interface{})
	// Slowv logs a message at slow level.
	Slowv(interface{})
	// Sloww logs a message at slow level.
	Sloww(string, ...LogField)
	// WithCallerSkip returns a new logger with the given caller skip.
	WithCallerSkip(skip int) Logger
	// WithContext returns a new logger with the given context.
	WithContext(ctx context.Context) Logger
	// WithDuration returns a new logger with the given duration.
	WithDuration(d time.Duration) Logger
	// WithFields returns a new logger with the given fields.
	WithFields(fields ...LogField) Logger
}
