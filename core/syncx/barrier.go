package syncx

import "sync"

// A Barrier is used to facility the barrier on a resource.
type Barrier struct {
	lock sync.Mutex
}

// Guard guards the given fn on the resource.
func (b *Barrier) Guard(fn func()) {
	Guard(&b.lock, fn)
}

// Guard guards the given fn with lock.
func Guard(lock sync.Locker, fn func()) {
	lock.Lock()
	defer lock.Unlock()
	fn()
}
