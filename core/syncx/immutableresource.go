package syncx

import (
	"sync"
	"time"

	"github.com/tal-tech/go-zero/core/timex"
)

const defaultRefreshInterval = time.Second

type (
	ImmutableResourceOption func(resource *ImmutableResource)

	ImmutableResource struct {
		fetch           func() (interface{}, error)
		resource        interface{}
		err             error
		lock            sync.RWMutex
		refreshInterval time.Duration
		lastTime        *AtomicDuration
	}
)

func NewImmutableResource(fn func() (interface{}, error), opts ...ImmutableResourceOption) *ImmutableResource {
	// cannot use executors.LessExecutor because of cycle imports
	ir := ImmutableResource{
		fetch:           fn,
		refreshInterval: defaultRefreshInterval,
		lastTime:        NewAtomicDuration(),
	}
	for _, opt := range opts {
		opt(&ir)
	}
	return &ir
}

func (ir *ImmutableResource) Get() (interface{}, error) {
	ir.lock.RLock()
	resource := ir.resource
	ir.lock.RUnlock()
	if resource != nil {
		return resource, nil
	}

	ir.maybeRefresh(func() {
		res, err := ir.fetch()
		ir.lock.Lock()
		if err != nil {
			ir.err = err
		} else {
			ir.resource, ir.err = res, nil
		}
		ir.lock.Unlock()
	})

	ir.lock.RLock()
	resource, err := ir.resource, ir.err
	ir.lock.RUnlock()
	return resource, err
}

func (ir *ImmutableResource) maybeRefresh(execute func()) {
	now := timex.Now()
	lastTime := ir.lastTime.Load()
	if lastTime == 0 || lastTime+ir.refreshInterval < now {
		ir.lastTime.Set(now)
		execute()
	}
}

// Set interval to 0 to enforce refresh every time if not succeeded. default is time.Second.
func WithRefreshIntervalOnFailure(interval time.Duration) ImmutableResourceOption {
	return func(resource *ImmutableResource) {
		resource.refreshInterval = interval
	}
}
