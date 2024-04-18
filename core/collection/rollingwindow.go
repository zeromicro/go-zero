package collection

import (
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/mathx"
	"github.com/zeromicro/go-zero/core/timex"
)

type (
	// BucketInterface is the interface that defines the buckets.
	BucketInterface[T Numerical] interface {
		Add(v T)
		Reset()
	}

	// Numerical is the interface that restricts the numerical type.
	Numerical = mathx.Numerical

	// RollingWindowOption let callers customize the RollingWindow.
	RollingWindowOption[T Numerical, B BucketInterface[T]] func(rollingWindow *RollingWindow[T, B])

	// RollingWindow defines a rolling window to calculate the events in buckets with the time interval.
	RollingWindow[T Numerical, B BucketInterface[T]] struct {
		lock          sync.RWMutex
		size          int
		win           *window[T, B]
		interval      time.Duration
		offset        int
		ignoreCurrent bool
		lastTime      time.Duration // start time of the last bucket
	}
)

// NewRollingWindow returns a RollingWindow that with size buckets and time interval,
// use opts to customize the RollingWindow.
func NewRollingWindow[T Numerical, B BucketInterface[T]](newBucket func() B, size int,
	interval time.Duration, opts ...RollingWindowOption[T, B]) *RollingWindow[T, B] {
	if size < 1 {
		panic("size must be greater than 0")
	}

	w := &RollingWindow[T, B]{
		size:     size,
		win:      newWindow[T, B](newBucket, size),
		interval: interval,
		lastTime: timex.Now(),
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Add adds value to current bucket.
func (rw *RollingWindow[T, B]) Add(v T) {
	rw.lock.Lock()
	defer rw.lock.Unlock()
	rw.updateOffset()
	rw.win.add(rw.offset, v)
}

// Reduce runs fn on all buckets, ignore current bucket if ignoreCurrent was set.
func (rw *RollingWindow[T, B]) Reduce(fn func(b B)) {
	rw.lock.RLock()
	defer rw.lock.RUnlock()

	var diff int
	span := rw.span()
	// ignore the current bucket, because of partial data
	if span == 0 && rw.ignoreCurrent {
		diff = rw.size - 1
	} else {
		diff = rw.size - span
	}
	if diff > 0 {
		offset := (rw.offset + span + 1) % rw.size
		rw.win.reduce(offset, diff, fn)
	}
}

func (rw *RollingWindow[T, B]) span() int {
	offset := int(timex.Since(rw.lastTime) / rw.interval)
	if 0 <= offset && offset < rw.size {
		return offset
	}

	return rw.size
}

func (rw *RollingWindow[T, B]) updateOffset() {
	span := rw.span()
	if span <= 0 {
		return
	}

	offset := rw.offset
	// reset expired buckets
	for i := 0; i < span; i++ {
		rw.win.resetBucket((offset + i + 1) % rw.size)
	}

	rw.offset = (offset + span) % rw.size
	now := timex.Now()
	// align to interval time boundary
	rw.lastTime = now - (now-rw.lastTime)%rw.interval
}

// Bucket defines the bucket that holds sum and num of additions.
type Bucket[T Numerical] struct {
	Sum   T
	Count int64
}

func (b *Bucket[T]) Add(v T) {
	b.Sum += v
	b.Count++
}

func (b *Bucket[T]) Reset() {
	b.Sum = 0
	b.Count = 0
}

type window[T Numerical, B BucketInterface[T]] struct {
	buckets []B
	size    int
}

func newWindow[T Numerical, B BucketInterface[T]](newBucket func() B, size int) *window[T, B] {
	buckets := make([]B, size)
	for i := 0; i < size; i++ {
		buckets[i] = newBucket()
	}
	return &window[T, B]{
		buckets: buckets,
		size:    size,
	}
}

func (w *window[T, B]) add(offset int, v T) {
	w.buckets[offset%w.size].Add(v)
}

func (w *window[T, B]) reduce(start, count int, fn func(b B)) {
	for i := 0; i < count; i++ {
		fn(w.buckets[(start+i)%w.size])
	}
}

func (w *window[T, B]) resetBucket(offset int) {
	w.buckets[offset%w.size].Reset()
}

// IgnoreCurrentBucket lets the Reduce call ignore current bucket.
func IgnoreCurrentBucket[T Numerical, B BucketInterface[T]]() RollingWindowOption[T, B] {
	return func(w *RollingWindow[T, B]) {
		w.ignoreCurrent = true
	}
}
