package syncx

import (
	"time"

	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/core/timex"
)

type Cond struct {
	signal chan lang.PlaceholderType
}

func NewCond() *Cond {
	return &Cond{
		signal: make(chan lang.PlaceholderType),
	}
}

// WaitWithTimeout wait for signal return remain wait time or timed out
func (cond *Cond) WaitWithTimeout(timeout time.Duration) (time.Duration, bool) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	begin := timex.Now()
	select {
	case <-cond.signal:
		elapsed := timex.Since(begin)
		remainTimeout := timeout - elapsed
		return remainTimeout, true
	case <-timer.C:
		return 0, false
	}
}

// Wait for signal
func (cond *Cond) Wait() {
	<-cond.signal
}

// Signal wakes one goroutine waiting on c, if there is any.
func (cond *Cond) Signal() {
	select {
	case cond.signal <- lang.Placeholder:
	default:
	}
}
