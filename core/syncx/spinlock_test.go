package syncx

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		lock.Lock()
		lock.Unlock()
		wait.Done()
	}()
	time.Sleep(time.Millisecond * 100)
	lock.Unlock()
	wait.Wait()
	assert.True(t, lock.TryLock())
}
