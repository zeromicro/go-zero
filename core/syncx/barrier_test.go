package syncx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBarrier_Guard(t *testing.T) {
	const total = 10000
	var barrier Barrier
	var count int
	for i := 0; i < total; i++ {
		barrier.Guard(func() {
			count++
		})
	}
	assert.Equal(t, total, count)
}

func TestBarrierPtr_Guard(t *testing.T) {
	const total = 10000
	barrier := new(Barrier)
	var count int
	for i := 0; i < total; i++ {
		barrier.Guard(func() {
			count++
		})
	}
	assert.Equal(t, total, count)
}
