package logx

import (
	"context"
	"time"
)

// A Logger represents a logger.
type Logger interface {
	Error(...interface{})
	Errorf(string, ...interface{})
	Errorv(interface{})
	Errorw(string, ...LogField)
	Info(...interface{})
	Infof(string, ...interface{})
	Infov(interface{})
	Infow(string, ...LogField)
	Slow(...interface{})
	Slowf(string, ...interface{})
	Slowv(interface{})
	Sloww(string, ...LogField)
	WithContext(context.Context) Logger
	WithDuration(time.Duration) Logger
}
