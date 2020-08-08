package collection

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stringx"
)

const duration = time.Millisecond * 50

func TestRollingWindowAdd(t *testing.T) {
	const size = 3
	r := NewRollingWindow(size, duration)
	listBuckets := func() []float64 {
		var buckets []float64
		r.Reduce(func(b *Bucket) {
			buckets = append(buckets, b.Sum)
		})
		return buckets
	}
	assert.Equal(t, []float64{0, 0, 0}, listBuckets())
	r.Add(1)
	assert.Equal(t, []float64{0, 0, 1}, listBuckets())
	elapse()
	r.Add(2)
	r.Add(3)
	assert.Equal(t, []float64{0, 1, 5}, listBuckets())
	elapse()
	r.Add(4)
	r.Add(5)
	r.Add(6)
	assert.Equal(t, []float64{1, 5, 15}, listBuckets())
	elapse()
	r.Add(7)
	assert.Equal(t, []float64{5, 15, 7}, listBuckets())
}

func TestRollingWindowReset(t *testing.T) {
	const size = 3
	r := NewRollingWindow(size, duration, IgnoreCurrentBucket())
	listBuckets := func() []float64 {
		var buckets []float64
		r.Reduce(func(b *Bucket) {
			buckets = append(buckets, b.Sum)
		})
		return buckets
	}
	r.Add(1)
	elapse()
	assert.Equal(t, []float64{0, 1}, listBuckets())
	elapse()
	assert.Equal(t, []float64{1}, listBuckets())
	elapse()
	assert.Nil(t, listBuckets())

	// cross window
	r.Add(1)
	time.Sleep(duration * 10)
	assert.Nil(t, listBuckets())
}

func TestRollingWindowReduce(t *testing.T) {
	const size = 4
	tests := []struct {
		win    *RollingWindow
		expect float64
	}{
		{
			win:    NewRollingWindow(size, duration),
			expect: 10,
		},
		{
			win:    NewRollingWindow(size, duration, IgnoreCurrentBucket()),
			expect: 4,
		},
	}

	for _, test := range tests {
		t.Run(stringx.Rand(), func(t *testing.T) {
			r := test.win
			for x := 0; x < size; x = x + 1 {
				for i := 0; i <= x; i++ {
					r.Add(float64(i))
				}
				if x < size-1 {
					elapse()
				}
			}
			var result float64
			r.Reduce(func(b *Bucket) {
				result += b.Sum
			})
			assert.Equal(t, test.expect, result)
		})
	}
}

func TestRollingWindowDataRace(t *testing.T) {
	const size = 3
	r := NewRollingWindow(size, duration)
	var stop = make(chan bool)
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				r.Add(float64(rand.Int63()))
				time.Sleep(duration / 2)
			}
		}
	}()
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				r.Reduce(func(b *Bucket) {})
			}
		}
	}()
	time.Sleep(duration * 5)
	close(stop)
}

func elapse() {
	time.Sleep(duration)
}
