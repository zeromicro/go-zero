package load

import (
	"errors"
	"fmt"
	"math"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mathx"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/timex"
)

const (
	defaultBuckets = 50
	defaultWindow  = time.Second * 5
	// using 1000m notation, 900m is like 90%, keep it as var for unit test
	defaultCpuThreshold = 900
	defaultMinRt        = float64(time.Second / time.Millisecond)
	// moving average hyperparameter beta for calculating requests on the fly
	flyingBeta               = 0.9
	coolOffDuration          = time.Second
	cpuMax                   = 1000 // millicpu
	millisecondsPerSecond    = 1000
	overloadFactorLowerBound = 0.1
)

var (
	// ErrServiceOverloaded is returned by Shedder.Allow when the service is overloaded.
	ErrServiceOverloaded = errors.New("service overloaded")

	// default to be enabled
	enabled = syncx.ForAtomicBool(true)
	// default to be enabled
	logEnabled = syncx.ForAtomicBool(true)
	// make it a variable for unit test
	systemOverloadChecker = func(cpuThreshold int64) bool {
		return stat.CpuUsage() >= cpuThreshold
	}
)

type (
	// A Promise interface is returned by Shedder.Allow to let callers tell
	// whether the processing request is successful or not.
	Promise interface {
		// Pass lets the caller tell that the call is successful.
		Pass()
		// Fail lets the caller tell that the call is failed.
		Fail()
	}

	// Shedder is the interface that wraps the Allow method.
	Shedder interface {
		// Allow returns the Promise if allowed, otherwise ErrServiceOverloaded.
		Allow() (Promise, error)
	}

	// ShedderOption lets caller customize the Shedder.
	ShedderOption func(opts *shedderOptions)

	shedderOptions struct {
		window       time.Duration
		buckets      int
		cpuThreshold int64
	}

	adaptiveShedder struct {
		cpuThreshold    int64
		windowScale     float64
		flying          int64
		avgFlying       float64
		avgFlyingLock   syncx.SpinLock
		overloadTime    *syncx.AtomicDuration
		droppedRecently *syncx.AtomicBool
		passCounter     *collection.RollingWindow[int64, *collection.Bucket[int64]]
		rtCounter       *collection.RollingWindow[int64, *collection.Bucket[int64]]
	}
)

// Disable lets callers disable load shedding.
func Disable() {
	enabled.Set(false)
}

// DisableLog disables the stat logs for load shedding.
func DisableLog() {
	logEnabled.Set(false)
}

// NewAdaptiveShedder returns an adaptive shedder.
// opts can be used to customize the Shedder.
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
	newBucket := func() *collection.Bucket[int64] {
		return new(collection.Bucket[int64])
	}
	return &adaptiveShedder{
		cpuThreshold:    options.cpuThreshold,
		windowScale:     float64(time.Second) / float64(bucketDuration) / millisecondsPerSecond,
		overloadTime:    syncx.NewAtomicDuration(),
		droppedRecently: syncx.NewAtomicBool(),
		passCounter:     collection.NewRollingWindow[int64, *collection.Bucket[int64]](newBucket, options.buckets, bucketDuration, collection.IgnoreCurrentBucket[int64, *collection.Bucket[int64]]()),
		rtCounter:       collection.NewRollingWindow[int64, *collection.Bucket[int64]](newBucket, options.buckets, bucketDuration, collection.IgnoreCurrentBucket[int64, *collection.Bucket[int64]]()),
	}
}

// Allow implements Shedder.Allow.
func (as *adaptiveShedder) Allow() (Promise, error) {
	if as.shouldDrop() {
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
	// this strategy makes avgFlying have a little bit of lag against flying, and smoother.
	// when the flying requests increase rapidly, avgFlying increase slower, accept more requests.
	// when the flying requests drop rapidly, avgFlying drop slower, accept fewer requests.
	// it makes the service to serve as many requests as possible.
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
	maxFlight := as.maxFlight() * as.overloadFactor()
	return avgFlying > maxFlight && float64(atomic.LoadInt64(&as.flying)) > maxFlight
}

func (as *adaptiveShedder) maxFlight() float64 {
	// windows = buckets per second
	// maxQPS = maxPASS * windows
	// minRT = min average response time in milliseconds
	// allowedFlying = maxQPS * minRT / milliseconds_per_second
	maxFlight := float64(as.maxPass()) * as.minRt() * as.windowScale
	return mathx.AtLeast(maxFlight, 1)
}

func (as *adaptiveShedder) maxPass() int64 {
	var result int64 = 1

	as.passCounter.Reduce(func(b *collection.Bucket[int64]) {
		if b.Sum > result {
			result = b.Sum
		}
	})

	return result
}

func (as *adaptiveShedder) minRt() float64 {
	// if no requests in previous windows, return defaultMinRt,
	// its a reasonable large value to avoid dropping requests.
	result := defaultMinRt

	as.rtCounter.Reduce(func(b *collection.Bucket[int64]) {
		if b.Count <= 0 {
			return
		}

		avg := math.Round(float64(b.Sum) / float64(b.Count))
		if avg < result {
			result = avg
		}
	})

	return result
}

func (as *adaptiveShedder) overloadFactor() float64 {
	// as.cpuThreshold must be less than cpuMax
	factor := (cpuMax - float64(stat.CpuUsage())) / (cpuMax - float64(as.cpuThreshold))
	// at least accept 10% of acceptable requests, even cpu is highly overloaded.
	return mathx.Between(factor, overloadFactorLowerBound, 1)
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

	overloadTime := as.overloadTime.Load()
	if overloadTime == 0 {
		return false
	}

	if timex.Since(overloadTime) < coolOffDuration {
		return true
	}

	as.droppedRecently.Set(false)
	return false
}

func (as *adaptiveShedder) systemOverloaded() bool {
	if !systemOverloadChecker(as.cpuThreshold) {
		return false
	}

	as.overloadTime.Set(timex.Now())
	return true
}

// WithBuckets customizes the Shedder with the given number of buckets.
func WithBuckets(buckets int) ShedderOption {
	return func(opts *shedderOptions) {
		opts.buckets = buckets
	}
}

// WithCpuThreshold customizes the Shedder with the given cpu threshold.
func WithCpuThreshold(threshold int64) ShedderOption {
	return func(opts *shedderOptions) {
		opts.cpuThreshold = threshold
	}
}

// WithWindow customizes the Shedder with given
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
	p.shedder.rtCounter.Add(int64(math.Ceil(rt)))
	p.shedder.passCounter.Add(1)
}
