package fx

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParallel(t *testing.T) {
	var count int32
	Parallel(func() {
		time.Sleep(time.Millisecond * 100)
		atomic.AddInt32(&count, 1)
	}, func() {
		time.Sleep(time.Millisecond * 100)
		atomic.AddInt32(&count, 2)
	}, func() {
		time.Sleep(time.Millisecond * 100)
		atomic.AddInt32(&count, 3)
	})
	assert.Equal(t, int32(6), count)
}
