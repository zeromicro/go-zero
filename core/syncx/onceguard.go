package syncx

import "sync/atomic"

type OnceGuard struct {
	done uint32
}

func (og *OnceGuard) Taken() bool {
	return atomic.LoadUint32(&og.done) == 1
}

func (og *OnceGuard) Take() bool {
	return atomic.CompareAndSwapUint32(&og.done, 0, 1)
}
