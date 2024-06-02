package internal

import (
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/grpclog"
)

// because grpclog.errorLog is not exported, we need to define our own.
const errorLevel = 2

// A Logger is a rpc logger.
type Logger struct{}

func init() {
	grpclog.SetLoggerV2(new(Logger))
}

// Error logs the given args into error log.
func (l *Logger) Error(args ...any) {
	logx.Error(args...)
}

// Errorf logs the given args with format into error log.
func (l *Logger) Errorf(format string, args ...any) {
	logx.Errorf(format, args...)
}

// Errorln logs the given args into error log with newline.
func (l *Logger) Errorln(args ...any) {
	logx.Error(args...)
}

// Fatal logs the given args into error log.
func (l *Logger) Fatal(args ...any) {
	logx.Error(args...)
}

// Fatalf logs the given args with format into error log.
func (l *Logger) Fatalf(format string, args ...any) {
	logx.Errorf(format, args...)
}

// Fatalln logs args into error log with newline.
func (l *Logger) Fatalln(args ...any) {
	logx.Error(args...)
}

// Info ignores the grpc info logs.
func (l *Logger) Info(_ ...any) {
	// ignore builtin grpc info
}

// Infoln ignores the grpc info logs.
func (l *Logger) Infoln(_ ...any) {
	// ignore builtin grpc info
}

// Infof ignores the grpc info logs.
func (l *Logger) Infof(_ string, _ ...any) {
	// ignore builtin grpc info
}

// V checks if meet required log level.
func (l *Logger) V(v int) bool {
	return v >= errorLevel
}

// Warning ignores the grpc warning logs.
func (l *Logger) Warning(_ ...any) {
	// ignore builtin grpc warning
}

// Warningf ignores the grpc warning logs.
func (l *Logger) Warningf(_ string, _ ...any) {
	// ignore builtin grpc warning
}

// Warningln ignores the grpc warning logs.
func (l *Logger) Warningln(_ ...any) {
	// ignore builtin grpc warning
}
