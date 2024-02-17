package executors

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChunkExecutor(t *testing.T) {
	var values []int
	var lock sync.Mutex

	executor := NewChunkExecutor(func(items []any) {
		lock.Lock()
		values = append(values, len(items))
		lock.Unlock()
	}, WithChunkBytes(10), WithFlushInterval(time.Minute))

	for i := 0; i < 50; i++ {
		executor.Add(1, 1)
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

func TestChunkExecutorFlushInterval(t *testing.T) {
	const (
		caches = 10
		size   = 5
	)
	var wait sync.WaitGroup

	wait.Add(1)
	executor := NewChunkExecutor(func(items []any) {
		assert.Equal(t, size, len(items))
		wait.Done()
	}, WithChunkBytes(caches), WithFlushInterval(time.Millisecond*100))

	for i := 0; i < size; i++ {
		executor.Add(1, 1)
	}

	wait.Wait()
}

func TestChunkExecutorEmpty(t *testing.T) {
	executor := NewChunkExecutor(func(items []any) {
		assert.Fail(t, "should not called")
	}, WithChunkBytes(10), WithFlushInterval(time.Millisecond))
	time.Sleep(time.Millisecond * 100)
	executor.Wait()
}

func TestChunkExecutorFlush(t *testing.T) {
	const (
		caches = 10
		tasks  = 5
	)

	var wait sync.WaitGroup
	wait.Add(1)
	be := NewChunkExecutor(func(items []any) {
		assert.Equal(t, tasks, len(items))
		wait.Done()
	}, WithChunkBytes(caches), WithFlushInterval(time.Minute))
	for i := 0; i < tasks; i++ {
		be.Add(1, 1)
	}
	be.Flush()
	wait.Wait()
}

func BenchmarkChunkExecutor(b *testing.B) {
	b.ReportAllocs()

	be := NewChunkExecutor(func(tasks []any) {
		time.Sleep(time.Millisecond * time.Duration(len(tasks)))
	})
	for i := 0; i < b.N; i++ {
		time.Sleep(time.Microsecond * 200)
		be.Add(1, 1)
	}
	be.Flush()
}
