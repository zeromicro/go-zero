package syncx

import (
	"runtime"
	"sync/atomic"
)

type SpinLock struct {
	lock uint32
}

func (sl *SpinLock) Lock() {
	for !sl.TryLock() {
		runtime.Gosched()
	}
}

func (sl *SpinLock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&sl.lock, 0, 1)
}

func (sl *SpinLock) Unlock() {
	atomic.StoreUint32(&sl.lock, 0)
}
