package executors

import (
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/threading"
)

// A DelayExecutor delays a tasks on given delay interval.
type DelayExecutor struct {
	fn        func()
	delay     time.Duration
	triggered bool
	lock      sync.Mutex
}

// NewDelayExecutor returns a DelayExecutor with given fn and delay.
func NewDelayExecutor(fn func(), delay time.Duration) *DelayExecutor {
	return &DelayExecutor{
		fn:    fn,
		delay: delay,
	}
}

// Trigger triggers the task to be executed after given delay, safe to trigger more than once.
func (de *DelayExecutor) Trigger() {
	de.lock.Lock()
	defer de.lock.Unlock()

	if de.triggered {
		return
	}

	de.triggered = true
	threading.GoSafe(func() {
		timer := time.NewTimer(de.delay)
		defer timer.Stop()
		<-timer.C

		// set triggered to false before calling fn to ensure no triggers are missed.
		de.lock.Lock()
		de.triggered = false
		de.lock.Unlock()
		de.fn()
	})
}
