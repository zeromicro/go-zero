package syncx

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeoutCondWait(t *testing.T) {
	var wait sync.WaitGroup
	cond := NewCond()
	wait.Add(2)
	go func() {
		cond.Wait()
		wait.Done()
	}()
	time.Sleep(time.Duration(50) * time.Millisecond)
	go func() {
		cond.Signal()
		wait.Done()
	}()
	wait.Wait()
}

func TestTimeoutCondWaitTimeout(t *testing.T) {
	var wait sync.WaitGroup
	cond := NewCond()
	wait.Add(1)
	go func() {
		cond.WaitWithTimeout(time.Duration(500) * time.Millisecond)
		wait.Done()
	}()
	wait.Wait()
}

func TestTimeoutCondWaitTimeoutRemain(t *testing.T) {
	var wait sync.WaitGroup
	cond := NewCond()
	wait.Add(2)
	ch := make(chan time.Duration, 1)
	defer close(ch)
	timeout := time.Duration(2000) * time.Millisecond
	go func() {
		remainTimeout, _ := cond.WaitWithTimeout(timeout)
		ch <- remainTimeout
		wait.Done()
	}()
	sleep(200)
	go func() {
		cond.Signal()
		wait.Done()
	}()
	wait.Wait()
	remainTimeout := <-ch
	assert.True(t, remainTimeout < timeout, "expect remainTimeout %v < %v", remainTimeout, timeout)
	assert.True(t, remainTimeout >= time.Duration(200)*time.Millisecond,
		"expect remainTimeout %v >= 200 millisecond", remainTimeout)
}

func TestSignalNoWait(t *testing.T) {
	cond := NewCond()
	cond.Signal()
}

func sleep(millisecond int) {
	time.Sleep(time.Duration(millisecond) * time.Millisecond)
}
