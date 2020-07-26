package syncx

import "sync"

type Barrier struct {
	lock sync.Mutex
}

func (b *Barrier) Guard(fn func()) {
	b.lock.Lock()
	defer b.lock.Unlock()
	fn()
}
