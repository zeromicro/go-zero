package syncx

import (
	"runtime"
	"sync/atomic"
)

// A SpinLock is used as a lock a fast execution.
type SpinLock struct {
	lock uint32
}

const maxBackoff = 64

// Lock locks the SpinLock.
// The Design reference https://github.com/panjf2000/ants/blob/master/internal/spinlock.go
func (sl *SpinLock) Lock() {
	backoff := 1
	for !sl.TryLock() {
		// Leverage the exponential backoff algorithm, see https://en.wikipedia.org/wiki/Exponential_backoff.
		for i := 0; i < backoff; i++ {
			runtime.Gosched()
		}
		if backoff < maxBackoff {
			backoff <<= 1
		}
	}
}

// TryLock tries to lock the SpinLock.
func (sl *SpinLock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&sl.lock, 0, 1)
}

// Unlock unlocks the SpinLock.
func (sl *SpinLock) Unlock() {
	atomic.StoreUint32(&sl.lock, 0)
}
