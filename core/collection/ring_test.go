package collection

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRing(t *testing.T) {
	assert.Panics(t, func() {
		NewRing(0)
	})
}

func TestRingLess(t *testing.T) {
	ring := NewRing(5)
	for i := 0; i < 3; i++ {
		ring.Add(i)
	}
	elements := ring.Take()
	assert.ElementsMatch(t, []any{0, 1, 2}, elements)
}

func TestRingMore(t *testing.T) {
	ring := NewRing(5)
	for i := 0; i < 11; i++ {
		ring.Add(i)
	}
	elements := ring.Take()
	assert.ElementsMatch(t, []any{6, 7, 8, 9, 10}, elements)
}

func TestRingAdd(t *testing.T) {
	ring := NewRing(5051)
	wg := sync.WaitGroup{}
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 1; j <= i; j++ {
				ring.Add(i)
			}
		}(i)
	}
	wg.Wait()
	assert.Equal(t, 5050, len(ring.Take()))
}

func BenchmarkRingAdd(b *testing.B) {
	ring := NewRing(500)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				ring.Add(i)
			}
		}
	})
}
