package collection

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stringx"
)

const duration = time.Millisecond * 50

func TestNewRollingWindow(t *testing.T) {
	assert.NotNil(t, NewRollingWindow[int64, *Bucket[int64]](func() *Bucket[int64] {
		return new(Bucket[int64])
	}, 10, time.Second))
	assert.Panics(t, func() {
		NewRollingWindow[int64, *Bucket[int64]](func() *Bucket[int64] {
			return new(Bucket[int64])
		}, 0, time.Second)
	})
}

func TestRollingWindowAdd(t *testing.T) {
	const size = 3
	r := NewRollingWindow[float64, *Bucket[float64]](func() *Bucket[float64] {
		return new(Bucket[float64])
	}, size, duration)
	listBuckets := func() []float64 {
		var buckets []float64
		r.Reduce(func(b *Bucket[float64]) {
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
	r := NewRollingWindow[float64, *Bucket[float64]](func() *Bucket[float64] {
		return new(Bucket[float64])
	}, size, duration, IgnoreCurrentBucket[float64, *Bucket[float64]]())
	listBuckets := func() []float64 {
		var buckets []float64
		r.Reduce(func(b *Bucket[float64]) {
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
		win    *RollingWindow[float64, *Bucket[float64]]
		expect float64
	}{
		{
			win: NewRollingWindow[float64, *Bucket[float64]](func() *Bucket[float64] {
				return new(Bucket[float64])
			}, size, duration),
			expect: 10,
		},
		{
			win: NewRollingWindow[float64, *Bucket[float64]](func() *Bucket[float64] {
				return new(Bucket[float64])
			}, size, duration, IgnoreCurrentBucket[float64, *Bucket[float64]]()),
			expect: 4,
		},
	}

	for _, test := range tests {
		t.Run(stringx.Rand(), func(t *testing.T) {
			r := test.win
			for x := 0; x < size; x++ {
				for i := 0; i <= x; i++ {
					r.Add(float64(i))
				}
				if x < size-1 {
					elapse()
				}
			}
			var result float64
			r.Reduce(func(b *Bucket[float64]) {
				result += b.Sum
			})
			assert.Equal(t, test.expect, result)
		})
	}
}

func TestRollingWindowBucketTimeBoundary(t *testing.T) {
	const size = 3
	interval := time.Millisecond * 30
	r := NewRollingWindow[float64, *Bucket[float64]](func() *Bucket[float64] {
		return new(Bucket[float64])
	}, size, interval)
	listBuckets := func() []float64 {
		var buckets []float64
		r.Reduce(func(b *Bucket[float64]) {
			buckets = append(buckets, b.Sum)
		})
		return buckets
	}
	assert.Equal(t, []float64{0, 0, 0}, listBuckets())
	r.Add(1)
	assert.Equal(t, []float64{0, 0, 1}, listBuckets())
	time.Sleep(time.Millisecond * 45)
	r.Add(2)
	r.Add(3)
	assert.Equal(t, []float64{0, 1, 5}, listBuckets())
	// sleep time should be less than interval, and make the bucket change happen
	time.Sleep(time.Millisecond * 20)
	r.Add(4)
	r.Add(5)
	r.Add(6)
	assert.Equal(t, []float64{1, 5, 15}, listBuckets())
	time.Sleep(time.Millisecond * 100)
	r.Add(7)
	r.Add(8)
	r.Add(9)
	assert.Equal(t, []float64{0, 0, 24}, listBuckets())
}

func TestRollingWindowDataRace(t *testing.T) {
	const size = 3
	r := NewRollingWindow[float64, *Bucket[float64]](func() *Bucket[float64] {
		return new(Bucket[float64])
	}, size, duration)
	stop := make(chan bool)
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
				r.Reduce(func(b *Bucket[float64]) {})
			}
		}
	}()
	time.Sleep(duration * 5)
	close(stop)
}

func elapse() {
	time.Sleep(duration)
}
