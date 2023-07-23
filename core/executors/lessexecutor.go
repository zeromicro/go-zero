package executors

import (
	"time"

	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/timex"
)

// A LessExecutor is an executor to limit execution once within given time interval.
type LessExecutor struct {
	threshold time.Duration
	lastTime  *syncx.AtomicDuration
}

// NewLessExecutor returns a LessExecutor with given threshold as time interval.
func NewLessExecutor(threshold time.Duration) *LessExecutor {
	return &LessExecutor{
		threshold: threshold,
		lastTime:  syncx.NewAtomicDuration(),
	}
}

// DoOrDiscard executes or discards the task depends on if
// another task was executed within the time interval.
func (le *LessExecutor) DoOrDiscard(execute func()) bool {
	now := timex.Now()
	lastTime := le.lastTime.Load()
	if lastTime == 0 || lastTime+le.threshold < now {
		le.lastTime.Set(now)
		execute()
		return true
	}

	return false
}
