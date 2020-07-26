package executors

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDelayExecutor(t *testing.T) {
	var count int32
	ex := NewDelayExecutor(func() {
		atomic.AddInt32(&count, 1)
	}, time.Millisecond*10)
	for i := 0; i < 100; i++ {
		ex.Trigger()
	}
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, int32(1), atomic.LoadInt32(&count))
}
