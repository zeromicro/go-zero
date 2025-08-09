package executors

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/core/timex"
)

const idleRound = 10

type (
	// TaskContainer interface defines a type that can be used as the underlying
	// container that used to do periodical executions.
	TaskContainer[T any] interface {
		// AddTask adds the task into the container.
		// Returns true if the container needs to be flushed after the addition.
		AddTask(task T) bool
		// Execute handles the collected tasks by the container when flushing.
		Execute(tasks []T)
		// RemoveAll removes the contained tasks, and return them.
		RemoveAll() []T
	}

	// A PeriodicalExecutor is an executor that periodically execute tasks.
	PeriodicalExecutor[T any] struct {
		commander chan []T
		interval  time.Duration
		container TaskContainer[T]
		waitGroup sync.WaitGroup
		// avoid race condition on waitGroup when calling wg.Add/Done/Wait(...)
		wgBarrier   syncx.Barrier
		confirmChan chan lang.PlaceholderType
		inflight    int32
		guarded     bool
		newTicker   func(duration time.Duration) timex.Ticker
		lock        sync.Mutex
	}
)

// NewPeriodicalExecutor returns a PeriodicalExecutor with given interval and container.
func NewPeriodicalExecutor[T any](interval time.Duration, container TaskContainer[T]) *PeriodicalExecutor[T] {
	executor := &PeriodicalExecutor[T]{
		// buffer 1 to let the caller go quickly
		commander:   make(chan []T, 1),
		interval:    interval,
		container:   container,
		confirmChan: make(chan lang.PlaceholderType),
		newTicker: func(d time.Duration) timex.Ticker {
			return timex.NewTicker(d)
		},
	}
	proc.AddShutdownListener(func() {
		executor.Flush()
	})

	return executor
}

// Add adds tasks into pe.
func (pe *PeriodicalExecutor[T]) Add(task T) {
	if vals, ok := pe.addAndCheck(task); ok {
		pe.commander <- vals
		<-pe.confirmChan
	}
}

// Flush forces pe to execute tasks.
func (pe *PeriodicalExecutor[T]) Flush() bool {
	pe.enterExecution()
	return pe.executeTasks(func() []T {
		pe.lock.Lock()
		defer pe.lock.Unlock()
		return pe.container.RemoveAll()
	}())
}

// Sync lets caller run fn thread-safe with pe, especially for the underlying container.
func (pe *PeriodicalExecutor[T]) Sync(fn func()) {
	pe.lock.Lock()
	defer pe.lock.Unlock()
	fn()
}

// Wait waits the execution to be done.
func (pe *PeriodicalExecutor[T]) Wait() {
	pe.Flush()
	pe.wgBarrier.Guard(func() {
		pe.waitGroup.Wait()
	})
}

func (pe *PeriodicalExecutor[T]) addAndCheck(task T) ([]T, bool) {
	pe.lock.Lock()
	defer func() {
		if !pe.guarded {
			pe.guarded = true
			// defer to unlock quickly
			defer pe.backgroundFlush()
		}
		pe.lock.Unlock()
	}()

	if pe.container.AddTask(task) {
		atomic.AddInt32(&pe.inflight, 1)
		return pe.container.RemoveAll(), true
	}

	return nil, false
}

func (pe *PeriodicalExecutor[T]) backgroundFlush() {
	go func() {
		// flush before quit goroutine to avoid missing tasks
		defer pe.Flush()

		ticker := pe.newTicker(pe.interval)
		defer ticker.Stop()

		var commanded bool
		last := timex.Now()
		for {
			select {
			case vals := <-pe.commander:
				commanded = true
				atomic.AddInt32(&pe.inflight, -1)
				pe.enterExecution()
				pe.confirmChan <- lang.Placeholder
				pe.executeTasks(vals)
				last = timex.Now()
			case <-ticker.Chan():
				if commanded {
					commanded = false
				} else if pe.Flush() {
					last = timex.Now()
				} else if pe.shallQuit(last) {
					return
				}
			}
		}
	}()
}

func (pe *PeriodicalExecutor[T]) doneExecution() {
	pe.waitGroup.Done()
}

func (pe *PeriodicalExecutor[T]) enterExecution() {
	pe.wgBarrier.Guard(func() {
		pe.waitGroup.Add(1)
	})
}

func (pe *PeriodicalExecutor[T]) executeTasks(tasks []T) bool {
	defer pe.doneExecution()

	ok := pe.hasTasks(tasks)
	if ok {
		threading.RunSafe(func() {
			pe.container.Execute(tasks)
		})
	}

	return ok
}

func (pe *PeriodicalExecutor[T]) hasTasks(tasks []T) bool {
	return len(tasks) > 0
}

func (pe *PeriodicalExecutor[T]) shallQuit(last time.Duration) (stop bool) {
	if timex.Since(last) <= pe.interval*idleRound {
		return
	}

	// checking pe.inflight and setting pe.guarded should be locked together
	pe.lock.Lock()
	if atomic.LoadInt32(&pe.inflight) == 0 {
		pe.guarded = false
		stop = true
	}
	pe.lock.Unlock()

	return
}
