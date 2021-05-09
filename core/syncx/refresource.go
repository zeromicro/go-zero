package syncx

import (
	"errors"
	"sync"
)

// ErrUseOfCleaned is an error that indicates using a cleaned resource.
var ErrUseOfCleaned = errors.New("using a cleaned resource")

// A RefResource is used to reference counting a resource.
type RefResource struct {
	lock    sync.Mutex
	ref     int32
	cleaned bool
	clean   func()
}

// NewRefResource returns a RefResource.
func NewRefResource(clean func()) *RefResource {
	return &RefResource{
		clean: clean,
	}
}

// Use uses the resource with reference count incremented.
func (r *RefResource) Use() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.cleaned {
		return ErrUseOfCleaned
	}

	r.ref++
	return nil
}

// Clean cleans a resource with reference count decremented.
func (r *RefResource) Clean() {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.cleaned {
		return
	}

	r.ref--
	if r.ref == 0 {
		r.cleaned = true
		r.clean()
	}
}
