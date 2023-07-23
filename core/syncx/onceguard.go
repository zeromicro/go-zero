package syncx

import "sync/atomic"

// An OnceGuard is used to make sure a resource can be taken once.
type OnceGuard struct {
	done uint32
}

// Taken checks if the resource is taken.
func (og *OnceGuard) Taken() bool {
	return atomic.LoadUint32(&og.done) == 1
}

// Take takes the resource, returns true on success, false for otherwise.
func (og *OnceGuard) Take() bool {
	return atomic.CompareAndSwapUint32(&og.done, 0, 1)
}
