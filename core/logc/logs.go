package logc

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type (
	LogConf  = logx.LogConf
	LogField = logx.LogField
)

// AddGlobalFields adds global fields.
func AddGlobalFields(fields ...LogField) {
	logx.AddGlobalFields(fields...)
}

// Alert alerts v in alert level, and the message is written to error log.
func Alert(_ context.Context, v string) {
	logx.Alert(v)
}

// Close closes the logging.
func Close() error {
	return logx.Close()
}

// Debug writes v into access log.
func Debug(ctx context.Context, v ...interface{}) {
	getLogger(ctx).Debug(v...)
}

// Debugf writes v with format into access log.
func Debugf(ctx context.Context, format string, v ...interface{}) {
	getLogger(ctx).Debugf(format, v...)
}

// Debugfn writes fn result into access log.
// This is useful when the function is expensive to compute,
// and we want to log it only when necessary.
func Debugfn(ctx context.Context, fn func() any) {
	getLogger(ctx).Debugfn(fn)
}

// Debugv writes v into access log with json content.
func Debugv(ctx context.Context, v interface{}) {
	getLogger(ctx).Debugv(v)
}

// Debugw writes msg along with fields into the access log.
func Debugw(ctx context.Context, msg string, fields ...LogField) {
	getLogger(ctx).Debugw(msg, fields...)
}

// Error writes v into error log.
func Error(ctx context.Context, v ...any) {
	getLogger(ctx).Error(v...)
}

// Errorf writes v with format into error log.
func Errorf(ctx context.Context, format string, v ...any) {
	getLogger(ctx).Errorf(fmt.Errorf(format, v...).Error())
}

// Errorfn writes fn result into error log.
// This is useful when the function is expensive to compute,
// and we want to log it only when necessary.
func Errorfn(ctx context.Context, fn func() any) {
	getLogger(ctx).Errorfn(fn)
}

// Errorv writes v into error log with json content.
// No call stack attached, because not elegant to pack the messages.
func Errorv(ctx context.Context, v any) {
	getLogger(ctx).Errorv(v)
}

// Errorw writes msg along with fields into the error log.
func Errorw(ctx context.Context, msg string, fields ...LogField) {
	getLogger(ctx).Errorw(msg, fields...)
}

// Field returns a LogField for the given key and value.
func Field(key string, value any) LogField {
	return logx.Field(key, value)
}

// Info writes v into access log.
func Info(ctx context.Context, v ...any) {
	getLogger(ctx).Info(v...)
}

// Infof writes v with format into access log.
func Infof(ctx context.Context, format string, v ...any) {
	getLogger(ctx).Infof(format, v...)
}

// Infofn writes fn result into access log.
// This is useful when the function is expensive to compute,
// and we want to log it only when necessary.
func Infofn(ctx context.Context, fn func() any) {
	getLogger(ctx).Infofn(fn)
}

// Infov writes v into access log with json content.
func Infov(ctx context.Context, v any) {
	getLogger(ctx).Infov(v)
}

// Infow writes msg along with fields into the access log.
func Infow(ctx context.Context, msg string, fields ...LogField) {
	getLogger(ctx).Infow(msg, fields...)
}

// Must checks if err is nil, otherwise logs the error and exits.
func Must(err error) {
	logx.Must(err)
}

// MustSetup sets up logging with given config c. It exits on error.
func MustSetup(c logx.LogConf) {
	logx.MustSetup(c)
}

// SetLevel sets the logging level. It can be used to suppress some logs.
func SetLevel(level uint32) {
	logx.SetLevel(level)
}

// SetUp sets up the logx.
// If already set up, return nil.
// We allow SetUp to be called multiple times, because, for example,
// we need to allow different service frameworks to initialize logx respectively.
// The same logic for SetUp
func SetUp(c LogConf) error {
	return logx.SetUp(c)
}

// Slow writes v into slow log.
func Slow(ctx context.Context, v ...any) {
	getLogger(ctx).Slow(v...)
}

// Slowf writes v with format into slow log.
func Slowf(ctx context.Context, format string, v ...any) {
	getLogger(ctx).Slowf(format, v...)
}

// Slowfn writes fn result into slow log.
// This is useful when the function is expensive to compute,
// and we want to log it only when necessary.
func Slowfn(ctx context.Context, fn func() any) {
	getLogger(ctx).Slowfn(fn)
}

// Slowv writes v into slow log with json content.
func Slowv(ctx context.Context, v any) {
	getLogger(ctx).Slowv(v)
}

// Sloww writes msg along with fields into slow log.
func Sloww(ctx context.Context, msg string, fields ...LogField) {
	getLogger(ctx).Sloww(msg, fields...)
}

// getLogger returns the logx.Logger with the given ctx and correct caller.
func getLogger(ctx context.Context) logx.Logger {
	return logx.WithContext(ctx).WithCallerSkip(1)
}
