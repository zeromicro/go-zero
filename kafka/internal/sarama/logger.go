package sarama

import (
	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/core/logx"
)

const logPrefix = "[sarama] "

type logBridge struct {
	innerLogger logx.Logger
}

func (l *logBridge) Print(v ...any) {
	vs := append([]any{logPrefix}, v...)
	l.innerLogger.Info(vs...)
}

func (l *logBridge) Printf(format string, v ...any) {
	l.innerLogger.Infof(logPrefix+format, v...)
}

func (l *logBridge) Println(v ...any) {
	vs := append([]any{logPrefix}, v...)
	l.innerLogger.Info(vs...)
}

func (l *logBridge) Debugf(format string, v ...any) {
	l.innerLogger.Debugf(logPrefix+format, v...)
}

func (l *logBridge) Errorf(format string, v ...any) {
	l.innerLogger.Errorf(logPrefix+format, v...)
}

func init() { //nolint
	logx.Info("init sarama logger")
	saramaLogger := &logBridge{logx.WithCallerSkip(1)}
	sarama.Logger = saramaLogger
	sarama.DebugLogger = saramaLogger
}
