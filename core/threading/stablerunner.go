package threading

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
)

const factor = 10

var (
	ErrRunnerClosed = errors.New("runner closed")

	bufSize = runtime.NumCPU() * factor
)

// StableRunner is a runner that guarantees messages are taken out with the pushed order.
// This runner is typically useful for Kafka consumers with parallel processing.
type StableRunner[I, O any] struct {
	handle        func(I) O
	consumedIndex uint64
	writtenIndex  uint64
	ring          []*struct {
		value O
		valid bool
		lock  *sync.Mutex
		cond  *sync.Cond
	}
	limitChan chan struct{}
	tokenChan chan struct{}
	runner    *TaskRunner
	closed    atomic.Bool
}

// NewStableRunner returns a new StableRunner with given message processor fn.
func NewStableRunner[I, O any](fn func(I) O) *StableRunner[I, O] {
	ring := make([]*struct {
		value O
		valid bool
		lock  *sync.Mutex
		cond  *sync.Cond
	}, bufSize)
	for i := 0; i < bufSize; i++ {
		lock := new(sync.Mutex)
		cond := sync.NewCond(lock)
		ring[i] = &struct {
			value O
			valid bool
			lock  *sync.Mutex
			cond  *sync.Cond
		}{
			lock: lock,
			cond: cond,
		}
	}

	return &StableRunner[I, O]{
		handle:    fn,
		ring:      ring,
		limitChan: make(chan struct{}, bufSize),
		tokenChan: make(chan struct{}, bufSize),
		runner:    NewTaskRunner(runtime.NumCPU()),
	}
}

// Get returns the next processed message in order.
// This method should be called in one goroutine.
func (r *StableRunner[I, O]) Get() (O, error) {
	if r.closed.Load() {
		if atomic.LoadUint64(&r.consumedIndex) == atomic.LoadUint64(&r.writtenIndex) {
			var o O
			return o, ErrRunnerClosed
		}
	}

	defer func() {
		atomic.AddUint64(&r.consumedIndex, 1)
		<-r.limitChan
	}()

	<-r.tokenChan
	index := atomic.LoadUint64(&r.consumedIndex)
	offset := index % uint64(bufSize)
	holder := r.ring[offset]
	holder.lock.Lock()
	defer holder.lock.Unlock()
	if !holder.valid {
		holder.cond.Wait()
	}
	holder.valid = false

	return holder.value, nil
}

// Push pushes the message v into the runner and to be processed concurrently,
// after processed, it will be cached to let caller take it in pushing order.
func (r *StableRunner[I, O]) Push(v I) error {
	if r.closed.Load() {
		return ErrRunnerClosed
	}

	r.limitChan <- struct{}{}
	index := atomic.AddUint64(&r.writtenIndex, 1)

	r.runner.Schedule(func() {
		defer func() {
			r.tokenChan <- struct{}{}
		}()

		o := r.handle(v)
		offset := (index - 1) % uint64(bufSize)
		holder := r.ring[offset]
		holder.lock.Lock()
		defer holder.lock.Unlock()
		holder.value = o
		holder.valid = true
		holder.cond.Signal()
	})

	return nil
}

// Wait waits all the messages to be processed and taken from inner buffer.
func (r *StableRunner[I, O]) Wait() {
	r.closed.Store(true)
	r.runner.Wait()
	for atomic.LoadUint64(&r.consumedIndex) != atomic.LoadUint64(&r.writtenIndex) {
		runtime.Gosched()
	}
}
