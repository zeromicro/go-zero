package errorx

import (
	"errors"
	"sync"
)

// BatchError is an error that can hold multiple errors.
type BatchError struct {
	errs []error
	lock sync.RWMutex
}

// Add adds one or more non-nil errors to the BatchError instance.
func (be *BatchError) Add(errs ...error) {
	be.lock.Lock()
	defer be.lock.Unlock()

	for _, err := range errs {
		if err != nil {
			be.errs = append(be.errs, err)
		}
	}
}

// Err returns an error that represents all accumulated errors.
// It returns nil if there are no errors.
func (be *BatchError) Err() error {
	be.lock.RLock()
	defer be.lock.RUnlock()

	// If there are no non-nil errors, errors.Join(...) returns nil.
	return errors.Join(be.errs...)
}

// NotNil checks if there is at least one error inside the BatchError.
func (be *BatchError) NotNil() bool {
	be.lock.RLock()
	defer be.lock.RUnlock()

	return len(be.errs) > 0
}
