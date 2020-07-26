package syncx

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtomicFloat64(t *testing.T) {
	f := ForAtomicFloat64(100)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 100; i++ {
				f.Add(1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, float64(600), f.Load())
}
