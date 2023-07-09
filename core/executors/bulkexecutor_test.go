package executors

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBulkExecutor(t *testing.T) {
	var values []int
	var lock sync.Mutex

	executor := NewBulkExecutor(func(items []any) {
		lock.Lock()
		values = append(values, len(items))
		lock.Unlock()
	}, WithBulkTasks(10), WithBulkInterval(time.Minute))

	for i := 0; i < 50; i++ {
		executor.Add(1)
		time.Sleep(time.Millisecond)
	}

	lock.Lock()
	assert.True(t, len(values) > 0)
	// ignore last value
	for i := 0; i < len(values); i++ {
		assert.Equal(t, 10, values[i])
	}
	lock.Unlock()
}

func TestBulkExecutorFlushInterval(t *testing.T) {
	const (
		caches = 10
		size   = 5
	)
	var wait sync.WaitGroup

	wait.Add(1)
	executor := NewBulkExecutor(func(items []any) {
		assert.Equal(t, size, len(items))
		wait.Done()
	}, WithBulkTasks(caches), WithBulkInterval(time.Millisecond*100))

	for i := 0; i < size; i++ {
		executor.Add(1)
	}

	wait.Wait()
}

func TestBulkExecutorEmpty(t *testing.T) {
	NewBulkExecutor(func(items []any) {
		assert.Fail(t, "should not called")
	}, WithBulkTasks(10), WithBulkInterval(time.Millisecond))
	time.Sleep(time.Millisecond * 100)
}

func TestBulkExecutorFlush(t *testing.T) {
	const (
		caches = 10
		tasks  = 5
	)

	var wait sync.WaitGroup
	wait.Add(1)
	be := NewBulkExecutor(func(items []any) {
		assert.Equal(t, tasks, len(items))
		wait.Done()
	}, WithBulkTasks(caches), WithBulkInterval(time.Minute))
	for i := 0; i < tasks; i++ {
		be.Add(1)
	}
	be.Flush()
	wait.Wait()
}

func TestBulkExecutorFlushSlowTasks(t *testing.T) {
	const total = 1500
	lock := new(sync.Mutex)
	result := make([]any, 0, 10000)
	exec := NewBulkExecutor(func(tasks []any) {
		time.Sleep(time.Millisecond * 100)
		lock.Lock()
		defer lock.Unlock()
		result = append(result, tasks...)
	}, WithBulkTasks(1000))
	for i := 0; i < total; i++ {
		assert.Nil(t, exec.Add(i))
	}

	exec.Flush()
	exec.Wait()
	assert.Equal(t, total, len(result))
}

func BenchmarkBulkExecutor(b *testing.B) {
	b.ReportAllocs()

	be := NewBulkExecutor(func(tasks []any) {
		time.Sleep(time.Millisecond * time.Duration(len(tasks)))
	})
	for i := 0; i < b.N; i++ {
		time.Sleep(time.Microsecond * 200)
		be.Add(1)
	}
	be.Flush()
}
