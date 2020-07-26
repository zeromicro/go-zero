package logx

type LessLogger struct {
	*limitedExecutor
}

func NewLessLogger(milliseconds int) *LessLogger {
	return &LessLogger{
		limitedExecutor: newLimitedExecutor(milliseconds),
	}
}

func (logger *LessLogger) Error(v ...interface{}) {
	logger.logOrDiscard(func() {
		Error(v...)
	})
}

func (logger *LessLogger) Errorf(format string, v ...interface{}) {
	logger.logOrDiscard(func() {
		Errorf(format, v...)
	})
}
