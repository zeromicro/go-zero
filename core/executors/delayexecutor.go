package executors

import (
	"sync"
	"time"

	"github.com/tal-tech/go-zero/core/threading"
)

type DelayExecutor struct {
	fn        func()
	delay     time.Duration
	triggered bool
	lock      sync.Mutex
}

func NewDelayExecutor(fn func(), delay time.Duration) *DelayExecutor {
	return &DelayExecutor{
		fn:    fn,
		delay: delay,
	}
}

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
