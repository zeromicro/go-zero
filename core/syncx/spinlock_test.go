package syncx

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/lang"
)

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
