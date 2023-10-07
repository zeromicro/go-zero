package executors

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/timex"
)

const threshold = 10

type container struct {
	interval time.Duration
	tasks    []int
	execute  func(tasks any)
}

func newContainer(interval time.Duration, execute func(tasks any)) *container {
	return &container{
		interval: interval,
		execute:  execute,
	}
}

func (c *container) AddTask(task any) bool {
	c.tasks = append(c.tasks, task.(int))
	return len(c.tasks) > threshold
}

func (c *container) Execute(tasks any) {
	if c.execute != nil {
		c.execute(tasks)
	} else {
		time.Sleep(c.interval)
	}
}

func (c *container) RemoveAll() any {
	tasks := c.tasks
	c.tasks = nil
	return tasks
}

func TestPeriodicalExecutor_Sync(t *testing.T) {
	var done int32
	exec := NewPeriodicalExecutor(time.Second, newContainer(time.Millisecond*500, nil))
	exec.Sync(func() {
		atomic.AddInt32(&done, 1)
	})
	assert.Equal(t, int32(1), atomic.LoadInt32(&done))
}

func TestPeriodicalExecutor_QuitGoroutine(t *testing.T) {
	ticker := timex.NewFakeTicker()
	exec := NewPeriodicalExecutor(time.Millisecond, newContainer(time.Millisecond, nil))
	exec.newTicker = func(d time.Duration) timex.Ticker {
		return ticker
	}
	routines := runtime.NumGoroutine()
	exec.Add(1)
	ticker.Tick()
	ticker.Wait(time.Millisecond * idleRound * 2)
	ticker.Tick()
	ticker.Wait(time.Millisecond * idleRound)
	assert.Equal(t, routines, runtime.NumGoroutine())
	proc.Shutdown()
}

func TestPeriodicalExecutor_Bulk(t *testing.T) {
	ticker := timex.NewFakeTicker()
	var vals []int
	// avoid data race
	var lock sync.Mutex
	exec := NewPeriodicalExecutor(time.Millisecond, newContainer(time.Millisecond, func(tasks any) {
		t := tasks.([]int)
		for _, each := range t {
			lock.Lock()
			vals = append(vals, each)
			lock.Unlock()
		}
	}))
	exec.newTicker = func(d time.Duration) timex.Ticker {
		return ticker
	}
	for i := 0; i < threshold*10; i++ {
		if i%threshold == 5 {
			time.Sleep(time.Millisecond * idleRound * 2)
		}
		exec.Add(i)
	}
	ticker.Tick()
	ticker.Wait(time.Millisecond * idleRound * 2)
	ticker.Tick()
	ticker.Tick()
	ticker.Wait(time.Millisecond * idleRound)
	var expect []int
	for i := 0; i < threshold*10; i++ {
		expect = append(expect, i)
	}

	lock.Lock()
	assert.EqualValues(t, expect, vals)
	lock.Unlock()
}

func TestPeriodicalExecutor_Panic(t *testing.T) {
	// avoid data race
	var lock sync.Mutex
	ticker := timex.NewFakeTicker()

	var (
		executedTasks []int
		expected      []int
	)
	executor := NewPeriodicalExecutor(time.Millisecond, newContainer(time.Millisecond, func(tasks any) {
		tt := tasks.([]int)
		lock.Lock()
		executedTasks = append(executedTasks, tt...)
		lock.Unlock()
		if tt[0] == 0 {
			panic("test")
		}
	}))
	executor.newTicker = func(duration time.Duration) timex.Ticker {
		return ticker
	}
	for i := 0; i < 30; i++ {
		executor.Add(i)
		expected = append(expected, i)
	}
	ticker.Tick()
	ticker.Tick()
	time.Sleep(time.Millisecond)
	lock.Lock()
	assert.Equal(t, expected, executedTasks)
	lock.Unlock()
}

func TestPeriodicalExecutor_FlushPanic(t *testing.T) {
	var (
		executedTasks []int
		expected      []int
		lock          sync.Mutex
	)
	executor := NewPeriodicalExecutor(time.Millisecond, newContainer(time.Millisecond, func(tasks any) {
		tt := tasks.([]int)
		lock.Lock()
		executedTasks = append(executedTasks, tt...)
		lock.Unlock()
		if tt[0] == 0 {
			panic("flush panic")
		}
	}))
	for i := 0; i < 8; i++ {
		executor.Add(i)
		expected = append(expected, i)
	}
	executor.Flush()
	lock.Lock()
	assert.Equal(t, expected, executedTasks)
	lock.Unlock()
}

func TestPeriodicalExecutor_Wait(t *testing.T) {
	var lock sync.Mutex
	executor := NewBulkExecutor(func(tasks []any) {
		lock.Lock()
		defer lock.Unlock()
		time.Sleep(10 * time.Millisecond)
	}, WithBulkTasks(1), WithBulkInterval(time.Second))
	for i := 0; i < 10; i++ {
		executor.Add(1)
	}
	executor.Flush()
	executor.Wait()
}

func TestPeriodicalExecutor_WaitFast(t *testing.T) {
	const total = 3
	var cnt int
	var lock sync.Mutex
	executor := NewBulkExecutor(func(tasks []any) {
		defer func() {
			cnt++
		}()
		lock.Lock()
		defer lock.Unlock()
		time.Sleep(10 * time.Millisecond)
	}, WithBulkTasks(1), WithBulkInterval(10*time.Millisecond))
	for i := 0; i < total; i++ {
		executor.Add(2)
	}
	executor.Flush()
	executor.Wait()
	assert.Equal(t, total, cnt)
}

func TestPeriodicalExecutor_Deadlock(t *testing.T) {
	executor := NewBulkExecutor(func(tasks []any) {
	}, WithBulkTasks(1), WithBulkInterval(time.Millisecond))
	for i := 0; i < 1e5; i++ {
		executor.Add(1)
	}
}

func TestPeriodicalExecutor_hasTasks(t *testing.T) {
	exec := NewPeriodicalExecutor(time.Millisecond, newContainer(time.Millisecond, nil))
	assert.False(t, exec.hasTasks(nil))
	assert.True(t, exec.hasTasks(1))
}

// go test -benchtime 10s -bench .
func BenchmarkExecutor(b *testing.B) {
	b.ReportAllocs()

	executor := NewPeriodicalExecutor(time.Second, newContainer(time.Millisecond*500, nil))
	for i := 0; i < b.N; i++ {
		executor.Add(1)
	}
}
