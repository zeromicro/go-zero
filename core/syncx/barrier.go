package syncx

import "sync"

// A Barrier is used to facility the barrier on a resource.
type Barrier struct {
	lock sync.Mutex
}

// Guard guards the given fn on the resource.
func (b *Barrier) Guard(fn func()) {
	b.lock.Lock()
	defer b.lock.Unlock()
	fn()
}
