package threading

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
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
		value chan O
		lock  sync.Mutex
	}
	runner *TaskRunner
	done   chan struct{}
}

// NewStableRunner returns a new StableRunner with given message processor fn.
func NewStableRunner[I, O any](fn func(I) O) *StableRunner[I, O] {
	ring := make([]*struct {
		value chan O
		lock  sync.Mutex
	}, bufSize)
	for i := 0; i < bufSize; i++ {
		ring[i] = &struct {
			value chan O
			lock  sync.Mutex
		}{
			value: make(chan O, 1),
		}
	}

	return &StableRunner[I, O]{
		handle: fn,
		ring:   ring,
		runner: NewTaskRunner(runtime.NumCPU()),
		done:   make(chan struct{}),
	}
}

// Get returns the next processed message in order.
// This method should be called in one goroutine.
func (r *StableRunner[I, O]) Get() (O, error) {
	defer atomic.AddUint64(&r.consumedIndex, 1)

	index := atomic.LoadUint64(&r.consumedIndex)
	offset := index % uint64(bufSize)
	holder := r.ring[offset]

	select {
	case o := <-holder.value:
		return o, nil
	case <-r.done:
		if atomic.LoadUint64(&r.consumedIndex) < atomic.LoadUint64(&r.writtenIndex) {
			return <-holder.value, nil
		}

		var o O
		return o, ErrRunnerClosed
	}
}

// Push pushes the message v into the runner and to be processed concurrently,
// after processed, it will be cached to let caller take it in pushing order.
func (r *StableRunner[I, O]) Push(v I) error {
	select {
	case <-r.done:
		return ErrRunnerClosed
	default:
		index := atomic.AddUint64(&r.writtenIndex, 1)
		offset := (index - 1) % uint64(bufSize)
		holder := r.ring[offset]
		holder.lock.Lock()
		r.runner.Schedule(func() {
			defer holder.lock.Unlock()
			o := r.handle(v)
			holder.value <- o
		})

		return nil
	}
}

// Wait waits all the messages to be processed and taken from inner buffer.
func (r *StableRunner[I, O]) Wait() {
	close(r.done)
	r.runner.Wait()
	for atomic.LoadUint64(&r.consumedIndex) < atomic.LoadUint64(&r.writtenIndex) {
		time.Sleep(time.Millisecond)
	}
}
