package syncx

import (
	"sync/atomic"
	"time"
)

// An AtomicDuration is an implementation of atomic duration.
type AtomicDuration atomic.Int64

// NewAtomicDuration returns an AtomicDuration.
func NewAtomicDuration() (d AtomicDuration) {
	d.Set(0)
	return
}

// ForAtomicDuration returns an AtomicDuration with given value.
func ForAtomicDuration(val time.Duration) (i AtomicDuration) {
	i.Set(val)
	return
}

// CompareAndSwap compares current value with old, if equals, set the value to val.
func (d *AtomicDuration) CompareAndSwap(old, val time.Duration) bool {
	return (*atomic.Int64)(d).CompareAndSwap(int64(old), int64(val))
}

// Load loads the current duration.
func (d *AtomicDuration) Load() time.Duration {
	return time.Duration((*atomic.Int64)(d).Load())
}

// Set sets the value to val.
func (d *AtomicDuration) Set(val time.Duration) {
	(*atomic.Int64)(d).Store(int64(val))
}
