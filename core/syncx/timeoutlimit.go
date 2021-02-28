package syncx

import (
	"errors"
	"time"
)

// ErrTimeout is an error that indicates the borrow timeout.
var ErrTimeout = errors.New("borrow timeout")

// A TimeoutLimit is used to borrow with timeouts.
type TimeoutLimit struct {
	limit Limit
	cond  *Cond
}

// NewTimeoutLimit returns a TimeoutLimit.
func NewTimeoutLimit(n int) TimeoutLimit {
	return TimeoutLimit{
		limit: NewLimit(n),
		cond:  NewCond(),
	}
}

// Borrow borrows with given timeout.
func (l TimeoutLimit) Borrow(timeout time.Duration) error {
	if l.TryBorrow() {
		return nil
	}

	var ok bool
	for {
		timeout, ok = l.cond.WaitWithTimeout(timeout)
		if ok && l.TryBorrow() {
			return nil
		}

		if timeout <= 0 {
			return ErrTimeout
		}
	}
}

// Return returns a borrow.
func (l TimeoutLimit) Return() error {
	if err := l.limit.Return(); err != nil {
		return err
	}

	l.cond.Signal()
	return nil
}

// TryBorrow tries a borrow.
func (l TimeoutLimit) TryBorrow() bool {
	return l.limit.TryBorrow()
}
