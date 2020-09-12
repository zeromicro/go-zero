package load

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/collection"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/mathx"
	"github.com/tal-tech/go-zero/core/stat"
	"github.com/tal-tech/go-zero/core/syncx"
)

const (
	buckets        = 10
	bucketDuration = time.Millisecond * 50
)

func init() {
	stat.SetReporter(nil)
}

func TestAdaptiveShedder(t *testing.T) {
	shedder := NewAdaptiveShedder(WithWindow(bucketDuration), WithBuckets(buckets), WithCpuThreshold(100))
	var wg sync.WaitGroup
	var drop int64
	proba := mathx.NewProba()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 30; i++ {
				promise, err := shedder.Allow()
				if err != nil {
					atomic.AddInt64(&drop, 1)
				} else {
					count := rand.Intn(5)
					time.Sleep(time.Millisecond * time.Duration(count))
					if proba.TrueOnProba(0.01) {
						promise.Fail()
					} else {
						promise.Pass()
					}
				}
			}
		}()
	}
	wg.Wait()
}

func TestAdaptiveShedderMaxPass(t *testing.T) {
	passCounter := newRollingWindow()
	for i := 1; i <= 10; i++ {
		passCounter.Add(float64(i * 100))
		time.Sleep(bucketDuration)
	}
	shedder := &adaptiveShedder{
		passCounter:     passCounter,
		droppedRecently: syncx.NewAtomicBool(),
	}
	assert.Equal(t, int64(1000), shedder.maxPass())

	// default max pass is equal to 1.
	passCounter = newRollingWindow()
	shedder = &adaptiveShedder{
		passCounter:     passCounter,
		droppedRecently: syncx.NewAtomicBool(),
	}
	assert.Equal(t, int64(1), shedder.maxPass())
}

func TestAdaptiveShedderMinRt(t *testing.T) {
	rtCounter := newRollingWindow()
	for i := 0; i < 10; i++ {
		if i > 0 {
			time.Sleep(bucketDuration)
		}
		for j := i*10 + 1; j <= i*10+10; j++ {
			rtCounter.Add(float64(j))
		}
	}
	shedder := &adaptiveShedder{
		rtCounter: rtCounter,
	}
	assert.Equal(t, float64(6), shedder.minRt())

	// default max min rt is equal to maxFloat64.
	rtCounter = newRollingWindow()
	shedder = &adaptiveShedder{
		rtCounter:       rtCounter,
		droppedRecently: syncx.NewAtomicBool(),
	}
	assert.Equal(t, defaultMinRt, shedder.minRt())
}

func TestAdaptiveShedderMaxFlight(t *testing.T) {
	passCounter := newRollingWindow()
	rtCounter := newRollingWindow()
	for i := 0; i < 10; i++ {
		if i > 0 {
			time.Sleep(bucketDuration)
		}
		passCounter.Add(float64((i + 1) * 100))
		for j := i*10 + 1; j <= i*10+10; j++ {
			rtCounter.Add(float64(j))
		}
	}
	shedder := &adaptiveShedder{
		passCounter:     passCounter,
		rtCounter:       rtCounter,
		windows:         buckets,
		droppedRecently: syncx.NewAtomicBool(),
	}
	assert.Equal(t, int64(54), shedder.maxFlight())
}

func TestAdaptiveShedderShouldDrop(t *testing.T) {
	logx.Disable()
	passCounter := newRollingWindow()
	rtCounter := newRollingWindow()
	for i := 0; i < 10; i++ {
		if i > 0 {
			time.Sleep(bucketDuration)
		}
		passCounter.Add(float64((i + 1) * 100))
		for j := i*10 + 1; j <= i*10+10; j++ {
			rtCounter.Add(float64(j))
		}
	}
	shedder := &adaptiveShedder{
		passCounter:     passCounter,
		rtCounter:       rtCounter,
		windows:         buckets,
		droppedRecently: syncx.NewAtomicBool(),
	}
	// cpu >=  800, inflight < maxPass
	systemOverloadChecker = func(int64) bool {
		return true
	}
	shedder.avgFlying = 50
	assert.False(t, shedder.shouldDrop())

	// cpu >=  800, inflight > maxPass
	shedder.avgFlying = 80
	shedder.flying = 50
	assert.False(t, shedder.shouldDrop())

	// cpu >=  800, inflight > maxPass
	shedder.avgFlying = 80
	shedder.flying = 80
	assert.True(t, shedder.shouldDrop())

	// cpu < 800, inflight > maxPass
	systemOverloadChecker = func(int64) bool {
		return false
	}
	shedder.avgFlying = 80
	assert.False(t, shedder.shouldDrop())
}

func BenchmarkAdaptiveShedder_Allow(b *testing.B) {
	logx.Disable()

	bench := func(b *testing.B) {
		var shedder = NewAdaptiveShedder()
		proba := mathx.NewProba()
		for i := 0; i < 6000; i++ {
			p, err := shedder.Allow()
			if err == nil {
				time.Sleep(time.Millisecond)
				if proba.TrueOnProba(0.01) {
					p.Fail()
				} else {
					p.Pass()
				}
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			p, err := shedder.Allow()
			if err == nil {
				p.Pass()
			}
		}
	}

	systemOverloadChecker = func(int64) bool {
		return true
	}
	b.Run("high load", bench)
	systemOverloadChecker = func(int64) bool {
		return false
	}
	b.Run("low load", bench)
}

func newRollingWindow() *collection.RollingWindow {
	return collection.NewRollingWindow(buckets, bucketDuration, collection.IgnoreCurrentBucket())
}
