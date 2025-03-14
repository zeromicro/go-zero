package logx

import (
	"context"
	"time"
)

// A Logger represents a logger.
type Logger interface {
	// Debug logs a message at debug level.
	Debug(...any)
	// Debugf logs a message at debug level.
	Debugf(string, ...any)
	// Debugfn logs a message at debug level.
	Debugfn(func() any)
	// Debugv logs a message at debug level.
	Debugv(any)
	// Debugw logs a message at debug level.
	Debugw(string, ...LogField)
	// Error logs a message at error level.
	Error(...any)
	// Errorf logs a message at error level.
	Errorf(string, ...any)
	// Errorfn logs a message at error level.
	Errorfn(func() any)
	// Errorv logs a message at error level.
	Errorv(any)
	// Errorw logs a message at error level.
	Errorw(string, ...LogField)
	// Info logs a message at info level.
	Info(...any)
	// Infof logs a message at info level.
	Infof(string, ...any)
	// Infofn logs a message at info level.
	Infofn(func() any)
	// Infov logs a message at info level.
	Infov(any)
	// Infow logs a message at info level.
	Infow(string, ...LogField)
	// Slow logs a message at slow level.
	Slow(...any)
	// Slowf logs a message at slow level.
	Slowf(string, ...any)
	// Slowfn logs a message at slow level.
	Slowfn(func() any)
	// Slowv logs a message at slow level.
	Slowv(any)
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
