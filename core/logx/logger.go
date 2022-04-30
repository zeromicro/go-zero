package logx

import (
	"context"
	"time"
)

// A Logger represents a logger.
type Logger interface {
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
	// WithContext returns a new logger with the given context.
	WithContext(context.Context) Logger
	// WithDuration returns a new logger with the given duration.
	WithDuration(time.Duration) Logger
}
