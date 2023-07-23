package syncx

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBarrier_Guard(t *testing.T) {
	const total = 10000
	var barrier Barrier
	var count int
	var wg sync.WaitGroup
	wg.Add(total)
	for i := 0; i < total; i++ {
		go barrier.Guard(func() {
			count++
			wg.Done()
		})
	}
	wg.Wait()
	assert.Equal(t, total, count)
}

func TestBarrierPtr_Guard(t *testing.T) {
	const total = 10000
	barrier := new(Barrier)
	var count int
	wg := new(sync.WaitGroup)
	wg.Add(total)
	for i := 0; i < total; i++ {
		go barrier.Guard(func() {
			count++
			wg.Done()
		})
	}
	wg.Wait()
	assert.Equal(t, total, count)
}

func TestGuard(t *testing.T) {
	const total = 10000
	var count int
	var lock sync.Mutex
	wg := new(sync.WaitGroup)
	wg.Add(total)
	for i := 0; i < total; i++ {
		go Guard(&lock, func() {
			count++
			wg.Done()
		})
	}
	wg.Wait()
	assert.Equal(t, total, count)
}
