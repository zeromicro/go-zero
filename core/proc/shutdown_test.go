//go:build linux || darwin || freebsd

package proc

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShutdown(t *testing.T) {
	t.Cleanup(restoreSettings)

	SetTimeToForceQuit(time.Hour)
	shutdownLock.Lock()
	assert.Equal(t, time.Hour, waitTime)
	shutdownLock.Unlock()

	var val int
	called := AddWrapUpListener(func() {
		val++
	})
	WrapUp()
	called()
	assert.Equal(t, 1, val)

	called = AddShutdownListener(func() {
		val += 2
	})
	Shutdown()
	called()
	assert.Equal(t, 3, val)
}

func TestShutdownWithMultipleServices(t *testing.T) {
	t.Cleanup(restoreSettings)

	SetTimeToForceQuit(time.Hour)
	shutdownLock.Lock()
	assert.Equal(t, time.Hour, waitTime)
	shutdownLock.Unlock()

	var val int32
	called1 := AddShutdownListener(func() {
		atomic.AddInt32(&val, 1)
	})
	called2 := AddShutdownListener(func() {
		atomic.AddInt32(&val, 2)
	})
	Shutdown()
	called1()
	called2()

	assert.Equal(t, int32(3), atomic.LoadInt32(&val))
}

func TestWrapUpWithMultipleServices(t *testing.T) {
	t.Cleanup(restoreSettings)

	SetTimeToForceQuit(time.Hour)
	shutdownLock.Lock()
	assert.Equal(t, time.Hour, waitTime)
	shutdownLock.Unlock()

	var val int32
	called1 := AddWrapUpListener(func() {
		atomic.AddInt32(&val, 1)
	})
	called2 := AddWrapUpListener(func() {
		atomic.AddInt32(&val, 2)
	})
	WrapUp()
	called1()
	called2()

	assert.Equal(t, int32(3), atomic.LoadInt32(&val))
}

func TestNotifyMoreThanOnce(t *testing.T) {
	t.Cleanup(restoreSettings)

	ch := make(chan struct{}, 1)

	go func() {
		var val int
		called := AddWrapUpListener(func() {
			val++
		})
		WrapUp()
		WrapUp()
		called()
		assert.Equal(t, 1, val)

		called = AddShutdownListener(func() {
			val += 2
		})
		Shutdown()
		Shutdown()
		called()
		assert.Equal(t, 3, val)
		ch <- struct{}{}
	}()

	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("timeout, check error logs")
	}
}

func TestSetup(t *testing.T) {
	t.Run("valid time", func(t *testing.T) {
		defer restoreSettings()

		Setup(ShutdownConf{
			WrapUpTime: time.Second * 2,
			WaitTime:   time.Second * 30,
		})

		shutdownLock.Lock()
		assert.Equal(t, time.Second*2, wrapUpTime)
		assert.Equal(t, time.Second*30, waitTime)
		shutdownLock.Unlock()
	})

	t.Run("valid time", func(t *testing.T) {
		defer restoreSettings()

		Setup(ShutdownConf{})

		shutdownLock.Lock()
		assert.Equal(t, defaultWrapUpTime, wrapUpTime)
		assert.Equal(t, defaultWaitTime, waitTime)
		shutdownLock.Unlock()
	})
}

func restoreSettings() {
	shutdownLock.Lock()
	defer shutdownLock.Unlock()

	wrapUpTime = defaultWrapUpTime
	waitTime = defaultWaitTime
}
