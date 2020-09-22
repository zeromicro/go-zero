package load

import (
	"errors"
	"fmt"
	"math"
	"sync/atomic"
	"time"

	"github.com/tal-tech/go-zero/core/collection"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stat"
	"github.com/tal-tech/go-zero/core/syncx"
	"github.com/tal-tech/go-zero/core/timex"
)

const (
	defaultBuckets = 50
	defaultWindow  = time.Second * 5
	// using 1000m notation, 900m is like 80%, keep it as var for unit test
	defaultCpuThreshold = 900
	defaultMinRt        = float64(time.Second / time.Millisecond)
	// moving average hyperparameter beta for calculating requests on the fly
	flyingBeta      = 0.9
	coolOffDuration = time.Second
)

var (
	ErrServiceOverloaded = errors.New("service overloaded")

	// default to be enabled
	enabled = syncx.ForAtomicBool(true)
	// make it a variable for unit test
	systemOverloadChecker = func(cpuThreshold int64) bool {
		return stat.CpuUsage() >= cpuThreshold
	}
)

type (
	Promise interface {
		Pass()
		Fail()
	}

	Shedder interface {
		Allow() (Promise, error)
	}

	ShedderOption func(opts *shedderOptions)

	shedderOptions struct {
		window       time.Duration
		buckets      int
		cpuThreshold int64
	}

	adaptiveShedder struct {
		cpuThreshold    int64
		windows         int64
		flying          int64
		avgFlying       float64
		avgFlyingLock   syncx.SpinLock
		dropTime        *syncx.AtomicDuration
		droppedRecently *syncx.AtomicBool
		passCounter     *collection.RollingWindow
		rtCounter       *collection.RollingWindow
	}
)

func Disable() {
	enabled.Set(false)
}

func NewAdaptiveShedder(opts ...ShedderOption) Shedder {
	if !enabled.True() {
		return newNopShedder()
	}

	options := shedderOptions{
		window:       defaultWindow,
		buckets:      defaultBuckets,
		cpuThreshold: defaultCpuThreshold,
	}
	for _, opt := range opts {
		opt(&options)
	}
	bucketDuration := options.window / time.Duration(options.buckets)
	return &adaptiveShedder{
		cpuThreshold:    options.cpuThreshold,
		windows:         int64(time.Second / bucketDuration),
		dropTime:        syncx.NewAtomicDuration(),
		droppedRecently: syncx.NewAtomicBool(),
		passCounter: collection.NewRollingWindow(options.buckets, bucketDuration,
			collection.IgnoreCurrentBucket()),
		rtCounter: collection.NewRollingWindow(options.buckets, bucketDuration,
			collection.IgnoreCurrentBucket()),
	}
}

func (as *adaptiveShedder) Allow() (Promise, error) {
	if as.shouldDrop() {
		as.dropTime.Set(timex.Now())
		as.droppedRecently.Set(true)

		return nil, ErrServiceOverloaded
	}

	as.addFlying(1)

	return &promise{
		start:   timex.Now(),
		shedder: as,
	}, nil
}

func (as *adaptiveShedder) addFlying(delta int64) {
	flying := atomic.AddInt64(&as.flying, delta)
	// update avgFlying when the request is finished.
	// this strategy makes avgFlying have a little bit lag against flying, and smoother.
	// when the flying requests increase rapidly, avgFlying increase slower, accept more requests.
	// when the flying requests drop rapidly, avgFlying drop slower, accept less requests.
	// it makes the service to serve as more requests as possible.
	if delta < 0 {
		as.avgFlyingLock.Lock()
		as.avgFlying = as.avgFlying*flyingBeta + float64(flying)*(1-flyingBeta)
		as.avgFlyingLock.Unlock()
	}
}

func (as *adaptiveShedder) highThru() bool {
	as.avgFlyingLock.Lock()
	avgFlying := as.avgFlying
	as.avgFlyingLock.Unlock()
	maxFlight := as.maxFlight()
	return int64(avgFlying) > maxFlight && atomic.LoadInt64(&as.flying) > maxFlight
}

func (as *adaptiveShedder) maxFlight() int64 {
	// windows = buckets per second
	// maxQPS = maxPASS * windows
	// minRT = min average response time in milliseconds
	// maxQPS * minRT / milliseconds_per_second
	return int64(math.Max(1, float64(as.maxPass()*as.windows)*(as.minRt()/1e3)))
}

func (as *adaptiveShedder) maxPass() int64 {
	var result float64 = 1

	as.passCounter.Reduce(func(b *collection.Bucket) {
		if b.Sum > result {
			result = b.Sum
		}
	})

	return int64(result)
}

func (as *adaptiveShedder) minRt() float64 {
	var result = defaultMinRt

	as.rtCounter.Reduce(func(b *collection.Bucket) {
		if b.Count <= 0 {
			return
		}

		avg := math.Round(b.Sum / float64(b.Count))
		if avg < result {
			result = avg
		}
	})

	return result
}

func (as *adaptiveShedder) shouldDrop() bool {
	if as.systemOverloaded() || as.stillHot() {
		if as.highThru() {
			flying := atomic.LoadInt64(&as.flying)
			as.avgFlyingLock.Lock()
			avgFlying := as.avgFlying
			as.avgFlyingLock.Unlock()
			msg := fmt.Sprintf(
				"dropreq, cpu: %d, maxPass: %d, minRt: %.2f, hot: %t, flying: %d, avgFlying: %.2f",
				stat.CpuUsage(), as.maxPass(), as.minRt(), as.stillHot(), flying, avgFlying)
			logx.Error(msg)
			stat.Report(msg)
			return true
		}
	}

	return false
}

func (as *adaptiveShedder) stillHot() bool {
	if !as.droppedRecently.True() {
		return false
	}

	dropTime := as.dropTime.Load()
	if dropTime == 0 {
		return false
	}

	hot := timex.Since(dropTime) < coolOffDuration
	if !hot {
		as.droppedRecently.Set(false)
	}

	return hot
}

func (as *adaptiveShedder) systemOverloaded() bool {
	return systemOverloadChecker(as.cpuThreshold)
}

func WithBuckets(buckets int) ShedderOption {
	return func(opts *shedderOptions) {
		opts.buckets = buckets
	}
}

func WithCpuThreshold(threshold int64) ShedderOption {
	return func(opts *shedderOptions) {
		opts.cpuThreshold = threshold
	}
}

func WithWindow(window time.Duration) ShedderOption {
	return func(opts *shedderOptions) {
		opts.window = window
	}
}

type promise struct {
	start   time.Duration
	shedder *adaptiveShedder
}

func (p *promise) Fail() {
	p.shedder.addFlying(-1)
}

func (p *promise) Pass() {
	rt := float64(timex.Since(p.start)) / float64(time.Millisecond)
	p.shedder.addFlying(-1)
	p.shedder.rtCounter.Add(math.Ceil(rt))
	p.shedder.passCounter.Add(1)
}
