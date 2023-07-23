package errorx

import "sync/atomic"

// AtomicError defines an atomic error.
type AtomicError struct {
	err atomic.Value // error
}

// Set sets the error.
func (ae *AtomicError) Set(err error) {
	if err != nil {
		ae.err.Store(err)
	}
}

// Load returns the error.
func (ae *AtomicError) Load() error {
	if v := ae.err.Load(); v != nil {
		return v.(error)
	}
	return nil
}
