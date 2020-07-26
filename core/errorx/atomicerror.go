package errorx

import "sync"

type AtomicError struct {
	err  error
	lock sync.Mutex
}

func (ae *AtomicError) Set(err error) {
	ae.lock.Lock()
	ae.err = err
	ae.lock.Unlock()
}

func (ae *AtomicError) Load() error {
	ae.lock.Lock()
	err := ae.err
	ae.lock.Unlock()
	return err
}
