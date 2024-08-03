package errorx

import (
	"errors"
	"sync"
)

type (
	// A BatchError is an error that can hold multiple errors.
	BatchError struct {
		errs []error
		lock sync.RWMutex
	}

	errorArray []error
)

// Add adds errs to be, nil errors are ignored.
func (be *BatchError) Add(errs ...error) {
	be.lock.Lock()
	defer be.lock.Unlock()

	for _, err := range errs {
		if err != nil {
			be.errs = append(be.errs, err)
		}
	}
}

// Err returns an error that represents all errors.
func (be *BatchError) Err() error {
	be.lock.RLock()
	defer be.lock.RUnlock()

	return errors.Join(be.errs...)
}

// NotNil checks if any error inside.
func (be *BatchError) NotNil() bool {
	be.lock.RLock()
	defer be.lock.RUnlock()

	return len(be.errs) > 0
}
