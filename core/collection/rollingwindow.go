package collection

import (
	"sync"
	"time"

	"github.com/tal-tech/go-zero/core/timex"
)

type (
	RollingWindowOption func(rollingWindow *RollingWindow)

	RollingWindow struct {
		lock          sync.RWMutex
		size          int
		win           *window
		interval      time.Duration
		offset        int
		ignoreCurrent bool
		lastTime      time.Duration
	}
)

func NewRollingWindow(size int, interval time.Duration, opts ...RollingWindowOption) *RollingWindow {
	w := &RollingWindow{
		size:     size,
		win:      newWindow(size),
		interval: interval,
		lastTime: timex.Now(),
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (rw *RollingWindow) Add(v float64) {
	rw.lock.Lock()
	defer rw.lock.Unlock()
	rw.updateOffset()
	rw.win.add(rw.offset, v)
}

func (rw *RollingWindow) Reduce(fn func(b *Bucket)) {
	rw.lock.RLock()
	defer rw.lock.RUnlock()

	var diff int
	span := rw.span()
	// ignore current bucket, because of partial data
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

func (rw *RollingWindow) span() int {
	offset := int(timex.Since(rw.lastTime) / rw.interval)
	if 0 <= offset && offset < rw.size {
		return offset
	} else {
		return rw.size
	}
}

func (rw *RollingWindow) updateOffset() {
	span := rw.span()
	if span > 0 {
		offset := rw.offset
		// reset expired buckets
		start := offset + 1
		steps := start + span
		var remainder int
		if steps > rw.size {
			remainder = steps - rw.size
			steps = rw.size
		}
		for i := start; i < steps; i++ {
			rw.win.resetBucket(i)
			offset = i
		}
		for i := 0; i < remainder; i++ {
			rw.win.resetBucket(i)
			offset = i
		}
		rw.offset = offset
		rw.lastTime = timex.Now()
	}
}

type Bucket struct {
	Sum   float64
	Count int64
}

func (b *Bucket) add(v float64) {
	b.Sum += v
	b.Count++
}

func (b *Bucket) reset() {
	b.Sum = 0
	b.Count = 0
}

type window struct {
	buckets []*Bucket
	size    int
}

func newWindow(size int) *window {
	var buckets []*Bucket
	for i := 0; i < size; i++ {
		buckets = append(buckets, new(Bucket))
	}
	return &window{
		buckets: buckets,
		size:    size,
	}
}

func (w *window) add(offset int, v float64) {
	w.buckets[offset%w.size].add(v)
}

func (w *window) reduce(start, count int, fn func(b *Bucket)) {
	for i := 0; i < count; i++ {
		fn(w.buckets[(start+i)%len(w.buckets)])
	}
}

func (w *window) resetBucket(offset int) {
	w.buckets[offset].reset()
}

func IgnoreCurrentBucket() RollingWindowOption {
	return func(w *RollingWindow) {
		w.ignoreCurrent = true
	}
}
