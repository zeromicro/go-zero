package syncx

import (
	"runtime"
	"sync/atomic"
)

// A SpinLock is used as a lock a fast execution.
type SpinLock struct {
	lock uint32
}

// Lock locks the SpinLock.
func (sl *SpinLock) Lock() {
	for !sl.TryLock() {
		runtime.Gosched()
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
