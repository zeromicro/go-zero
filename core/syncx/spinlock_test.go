package syncx

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/lang"
)

type NoBackOffSpinLock struct {
	lock uint32
}

func (sl *NoBackOffSpinLock) Lock() {
	for !sl.TryLock() {
		runtime.Gosched()
	}
}

func (sl *NoBackOffSpinLock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&sl.lock, 0, 1)
}

func (sl *NoBackOffSpinLock) Unlock() {
	atomic.StoreUint32(&sl.lock, 0)
}

func GetNoBackOffSpinLock() sync.Locker {
	return new(NoBackOffSpinLock)
}

type BackOffSpinLock struct {
	lock uint32
}

func (sl *BackOffSpinLock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&sl.lock, 0, 1)
}

func (sl *BackOffSpinLock) Lock() {
	wait := 1
	for !sl.TryLock() {
		for i := 0; i < wait; i++ {
			runtime.Gosched()
		}
		if wait < maxBackoff {
			wait <<= 1
		}
	}
}

func (sl *BackOffSpinLock) Unlock() {
	atomic.StoreUint32(&sl.lock, 0)
}

func GetBackOffSpinLock() sync.Locker {
	return new(BackOffSpinLock)
}

func TestTryLock(t *testing.T) {
	var lock SpinLock
	assert.True(t, lock.TryLock())
	assert.False(t, lock.TryLock())
	lock.Unlock()
	assert.True(t, lock.TryLock())
}

func TestSpinLock(t *testing.T) {
	var lock SpinLock
	lock.Lock()
	assert.False(t, lock.TryLock())
	lock.Unlock()
	assert.True(t, lock.TryLock())
}

func TestSpinLockRace(t *testing.T) {
	var lock SpinLock
	lock.Lock()
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		wait.Done()
	}()
	time.Sleep(time.Millisecond * 100)
	lock.Unlock()
	wait.Wait()
	assert.True(t, lock.TryLock())
}

func TestSpinLock_TryLock(t *testing.T) {
	var lock SpinLock
	var count int32
	var wait sync.WaitGroup
	wait.Add(2)
	sig := make(chan lang.PlaceholderType)

	go func() {
		lock.TryLock()
		sig <- lang.Placeholder
		atomic.AddInt32(&count, 1)
		runtime.Gosched()
		lock.Unlock()
		wait.Done()
	}()

	go func() {
		<-sig
		lock.Lock()
		atomic.AddInt32(&count, 1)
		lock.Unlock()
		wait.Done()
	}()

	wait.Wait()
	assert.Equal(t, int32(2), atomic.LoadInt32(&count))
}

func BenchmarkMutex(b *testing.B) {
	count := 1
	m := sync.Mutex{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Lock()
			count = count + 1
			m.Unlock()
		}
	})
}

func BenchmarkSpinLock(b *testing.B) {
	count := 1
	spin := GetNoBackOffSpinLock()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			spin.Lock()
			count = count + 1
			spin.Unlock()
		}
	})
}

func BenchmarkBackOffSpinLock(b *testing.B) {
	count := 1
	spin := GetBackOffSpinLock()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			spin.Lock()
			count = count + 1
			spin.Unlock()
		}
	})
}

/*
Benchmark result for three types of locks:
goos: windows
goarch: amd64
cpu: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
BenchmarkMutex
BenchmarkMutex-8                26066272                46.79 ns/op            0 B/op          0 allocs/op
BenchmarkSpinLock
BenchmarkSpinLock-8             52740297                23.30 ns/op            0 B/op          0 allocs/op
BenchmarkBackOffSpinLock
BenchmarkBackOffSpinLock-8      69369703                17.32 ns/op            0 B/op          0 allocs/op

goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
BenchmarkMutex
BenchmarkMutex-4                24296307                41.86 ns/op            0 B/op          0 allocs/op
BenchmarkSpinLock
BenchmarkSpinLock-4             47955232                29.39 ns/op            0 B/op          0 allocs/op
BenchmarkBackOffSpinLock
BenchmarkBackOffSpinLock-4      48560684                22.79 ns/op            0 B/op          0 allocs/op

*/
