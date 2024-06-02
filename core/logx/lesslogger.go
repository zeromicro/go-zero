package logx

// A LessLogger is a logger that controls to log once during the given duration.
type LessLogger struct {
	*limitedExecutor
}

// NewLessLogger returns a LessLogger.
func NewLessLogger(milliseconds int) *LessLogger {
	return &LessLogger{
		limitedExecutor: newLimitedExecutor(milliseconds),
	}
}

// Error logs v into error log or discard it if more than once in the given duration.
func (logger *LessLogger) Error(v ...any) {
	logger.logOrDiscard(func() {
		Error(v...)
	})
}

// Errorf logs v with format into error log or discard it if more than once in the given duration.
func (logger *LessLogger) Errorf(format string, v ...any) {
	logger.logOrDiscard(func() {
		Errorf(format, v...)
	})
}
