//go:build linux || darwin

package proc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShutdown(t *testing.T) {
	SetTimeToForceQuit(time.Hour)
	assert.Equal(t, time.Hour, delayTimeBeforeForceQuit)

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

func TestNotifyMoreThanOnce(t *testing.T) {
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
	Setup(ProcConf{
		WrapUpTime: time.Second * 2,
		WaitTime:   time.Second * 30,
	})
	assert.Equal(t, time.Second*2, wrapUpTime)
	assert.Equal(t, time.Second*30, delayTimeBeforeForceQuit)
}
