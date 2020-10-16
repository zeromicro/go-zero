package errorx

import "sync/atomic"

type AtomicError struct {
	err atomic.Value // error
}

func (ae *AtomicError) Set(err error) {
	ae.err.Store(err)
}

func (ae *AtomicError) Load() error {
	if v := ae.err.Load(); v != nil {
		return v.(error)
	}
	return nil
}
