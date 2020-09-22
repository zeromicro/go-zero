package internal

import (
	"sync"

	"github.com/tal-tech/go-zero/core/logx"
	"google.golang.org/grpc/grpclog"
)

// because grpclog.errorLog is not exported, we need to define our own.
const errorLevel = 2

var once sync.Once

type Logger struct{}

func InitLogger() {
	once.Do(func() {
		grpclog.SetLoggerV2(new(Logger))
	})
}

func (l *Logger) Error(args ...interface{}) {
	logx.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	logx.Errorf(format, args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	logx.Error(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	logx.Error(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	logx.Errorf(format, args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	logx.Error(args...)
}

func (l *Logger) Info(args ...interface{}) {
	// ignore builtin grpc info
}

func (l *Logger) Infoln(args ...interface{}) {
	// ignore builtin grpc info
}

func (l *Logger) Infof(format string, args ...interface{}) {
	// ignore builtin grpc info
}

func (l *Logger) V(v int) bool {
	return v >= errorLevel
}

func (l *Logger) Warning(args ...interface{}) {
	// ignore builtin grpc warning
}

func (l *Logger) Warningln(args ...interface{}) {
	// ignore builtin grpc warning
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	// ignore builtin grpc warning
}
