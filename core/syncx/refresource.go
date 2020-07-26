package syncx

import (
	"errors"
	"sync"
)

var ErrUseOfCleaned = errors.New("using a cleaned resource")

type RefResource struct {
	lock    sync.Mutex
	ref     int32
	cleaned bool
	clean   func()
}

func NewRefResource(clean func()) *RefResource {
	return &RefResource{
		clean: clean,
	}
}

func (r *RefResource) Use() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.cleaned {
		return ErrUseOfCleaned
	}

	r.ref++
	return nil
}

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
