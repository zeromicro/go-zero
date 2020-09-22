package executors

import (
	"time"

	"github.com/tal-tech/go-zero/core/syncx"
	"github.com/tal-tech/go-zero/core/timex"
)

type LessExecutor struct {
	threshold time.Duration
	lastTime  *syncx.AtomicDuration
}

func NewLessExecutor(threshold time.Duration) *LessExecutor {
	return &LessExecutor{
		threshold: threshold,
		lastTime:  syncx.NewAtomicDuration(),
	}
}

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
