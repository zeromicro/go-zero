//go:build linux || darwin

package proc

import (
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestDone(t *testing.T) {
	select {
	case <-Done():
		assert.Fail(t, "should run")
	default:
	}
	assert.NotNil(t, Done())
}

func TestSIGTERMShutdownSignal(t *testing.T) {
	p, err := os.FindProcess(os.Getpid())
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)

	err = p.Signal(syscall.SIGTERM)
	assert.Nil(t, err)

	<-Done()
}

func TestSIGINTShutdownSignal(t *testing.T) {
	p, err := os.FindProcess(os.Getpid())
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)

	err = p.Signal(syscall.SIGTERM)
	assert.Nil(t, err)

	<-Done()
}
