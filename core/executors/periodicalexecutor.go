package executors

import (
	"reflect"
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
	TaskContainer interface {
		// AddTask adds the task into the container.
		// Returns true if the container needs to be flushed after the addition.
		AddTask(task any) bool
		// Execute handles the collected tasks by the container when flushing.
		Execute(tasks any)
		// RemoveAll removes the contained tasks, and return them.
		RemoveAll() any
	}

	// A PeriodicalExecutor is an executor that periodically execute tasks.
	PeriodicalExecutor struct {
		commander chan any
		interval  time.Duration
		container TaskContainer
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
func NewPeriodicalExecutor(interval time.Duration, container TaskContainer) *PeriodicalExecutor {
	executor := &PeriodicalExecutor{
		// buffer 1 to let the caller go quickly
		commander:   make(chan any, 1),
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
func (pe *PeriodicalExecutor) Add(task any) {
	if vals, ok := pe.addAndCheck(task); ok {
		pe.commander <- vals
		<-pe.confirmChan
	}
}

// Flush forces pe to execute tasks.
func (pe *PeriodicalExecutor) Flush() bool {
	pe.enterExecution()
	return pe.executeTasks(func() any {
		pe.lock.Lock()
		defer pe.lock.Unlock()
		return pe.container.RemoveAll()
	}())
}

// Sync lets caller run fn thread-safe with pe, especially for the underlying container.
func (pe *PeriodicalExecutor) Sync(fn func()) {
	pe.lock.Lock()
	defer pe.lock.Unlock()
	fn()
}

// Wait waits the execution to be done.
func (pe *PeriodicalExecutor) Wait() {
	pe.Flush()
	pe.wgBarrier.Guard(func() {
		pe.waitGroup.Wait()
	})
}

func (pe *PeriodicalExecutor) addAndCheck(task any) (any, bool) {
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

func (pe *PeriodicalExecutor) backgroundFlush() {
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

func (pe *PeriodicalExecutor) doneExecution() {
	pe.waitGroup.Done()
}

func (pe *PeriodicalExecutor) enterExecution() {
	pe.wgBarrier.Guard(func() {
		pe.waitGroup.Add(1)
	})
}

func (pe *PeriodicalExecutor) executeTasks(tasks any) bool {
	defer pe.doneExecution()

	ok := pe.hasTasks(tasks)
	if ok {
		threading.RunSafe(func() {
			pe.container.Execute(tasks)
		})
	}

	return ok
}

func (pe *PeriodicalExecutor) hasTasks(tasks any) bool {
	if tasks == nil {
		return false
	}

	val := reflect.ValueOf(tasks)
	switch val.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return val.Len() > 0
	default:
		// unknown type, let caller execute it
		return true
	}
}

func (pe *PeriodicalExecutor) shallQuit(last time.Duration) (stop bool) {
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
