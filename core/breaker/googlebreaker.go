package breaker

import (
	"math"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/mathx"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/timex"
)

const (
	// 250ms for bucket duration
	window                      = time.Second * 10
	buckets                     = 40
	forcePassDuration           = time.Second
	k                           = 1.5
	minK                        = 1.1
	protection                  = 5
	latencyActivationMultiplier = 3
	latencyCeilingRatio         = 0.95
	latencyBaselineDecayBeta    = 0.25
	latencyBaselineRiseBeta     = 0.01
	latencyCurrentBeta          = 0.25
	latencyMaxDropRatio         = 0.3
)

// googleBreaker is a netflixBreaker pattern from google.
// see Client-Side Throttling section in https://landing.google.com/sre/sre-book/chapters/handling-overload/
type (
	googleBreaker struct {
		k                float64
		stat             *collection.RollingWindow[int64, *bucket]
		proba            *mathx.Proba
		lastPass         *syncx.AtomicDuration
		timeoutUs        int64
		noLoadLatencyUs  int64
		currentLatencyUs int64
	}

	windowResult struct {
		accepts        int64
		total          int64
		failingBuckets int64
		workingBuckets int64
	}
)

func newGoogleBreaker(timeout time.Duration) *googleBreaker {
	bucketDuration := time.Duration(int64(window) / int64(buckets))
	st := collection.NewRollingWindow(func() *bucket {
		return new(bucket)
	}, buckets, bucketDuration)
	return &googleBreaker{
		stat:      st,
		k:         k,
		proba:     mathx.NewProba(),
		lastPass:  syncx.NewAtomicDuration(),
		timeoutUs: timeout.Microseconds(),
	}
}

func (b *googleBreaker) accept() error {
	history := b.history()
	errorRatio := b.calcK(history)
	latencyRatio := b.calcLatencyRatio()
	dropRatio := math.Max(errorRatio, latencyRatio)
	if dropRatio <= 0 {
		return nil
	}

	lastPass := b.lastPass.Load()
	if lastPass > 0 && timex.Since(lastPass) > forcePassDuration {
		b.lastPass.Set(timex.Now())
		return nil
	}

	dropRatio *= float64(buckets-history.workingBuckets) / buckets

	if b.proba.TrueOnProba(dropRatio) {
		return ErrServiceUnavailable
	}

	b.lastPass.Set(timex.Now())

	return nil
}

func (b *googleBreaker) allow() (internalPromise, error) {
	if err := b.accept(); err != nil {
		b.markDrop()
		return nil, err
	}

	return googlePromise{
		b:     b,
		start: timex.Now(),
	}, nil
}

func (b *googleBreaker) calcK(history windowResult) float64 {
	w := b.k - (b.k-minK)*float64(history.failingBuckets)/buckets
	weightedAccepts := mathx.AtLeast(w, minK) * float64(history.accepts)
	// https://landing.google.com/sre/sre-book/chapters/handling-overload/#eq2101
	// for better performance, no need to care about the negative ratio
	return (float64(history.total-protection) - weightedAccepts) / float64(history.total+1)
}

func (b *googleBreaker) calcLatencyRatio() float64 {
	if b.timeoutUs <= 0 {
		return 0
	}

	noLoadLatency := atomic.LoadInt64(&b.noLoadLatencyUs)
	currentLatencyUs := atomic.LoadInt64(&b.currentLatencyUs)
	if noLoadLatency <= 0 || currentLatencyUs <= 0 {
		return 0
	}

	threshold := noLoadLatency * latencyActivationMultiplier
	ceiling := int64(float64(b.timeoutUs) * latencyCeilingRatio)
	if currentLatencyUs < threshold || ceiling <= threshold {
		return 0
	}

	ratio := float64(currentLatencyUs-threshold) / float64(ceiling-threshold)
	return math.Min(1.0, math.Max(0.0, ratio)) * latencyMaxDropRatio
}

func (b *googleBreaker) doReq(req func() error, fallback Fallback, acceptable Acceptable) error {
	if err := b.accept(); err != nil {
		b.markDrop()
		if fallback != nil {
			return fallback(err)
		}

		return err
	}

	var succ bool
	start := timex.Now()
	defer func() {
		// if req() panic, success is false, mark as failure
		if succ {
			b.markSuccess(timex.Since(start).Microseconds())
		} else {
			b.markFailure()
		}
	}()

	err := req()
	if acceptable(err) {
		succ = true
	}

	return err
}

func (b *googleBreaker) history() windowResult {
	var result windowResult

	b.stat.Reduce(func(b *bucket) {
		result.accepts += b.Success
		result.total += b.Sum
		if b.Failure > 0 {
			result.workingBuckets = 0
		} else if b.Success > 0 {
			result.workingBuckets++
		}
		if b.Success > 0 {
			result.failingBuckets = 0
		} else if b.Failure > 0 {
			result.failingBuckets++
		}
	})

	return result
}

func (b *googleBreaker) markDrop() {
	b.stat.Add(drop)
}

func (b *googleBreaker) markFailure() {
	b.stat.Add(fail)
}

func (b *googleBreaker) markSuccess(latencyUs int64) {
	b.stat.Add(success)
	if b.timeoutUs > 0 {
		b.updateLatency(latencyUs)
	}
}

func (b *googleBreaker) updateBaselineLatency(latencyUs int64) {
	noLoadLatency := atomic.LoadInt64(&b.noLoadLatencyUs)
	if noLoadLatency <= 0 {
		atomic.StoreInt64(&b.noLoadLatencyUs, latencyUs)
		return
	}

	var beta float64
	if latencyUs < noLoadLatency {
		// Fast decay when latency decreases
		beta = latencyBaselineDecayBeta
	} else {
		// Slow rise when latency increases
		beta = latencyBaselineRiseBeta
	}

	newBaseline := int64(beta*float64(latencyUs) + (1-beta)*float64(noLoadLatency))
	atomic.StoreInt64(&b.noLoadLatencyUs, newBaseline)
}

func (b *googleBreaker) updateCurrentLatency(latencyUs int64) {
	currentLatency := atomic.LoadInt64(&b.currentLatencyUs)
	if currentLatency <= 0 {
		atomic.StoreInt64(&b.currentLatencyUs, latencyUs)
		return
	}

	// Fast EMA to update current latency
	newCurrent := int64(latencyCurrentBeta*float64(latencyUs) + (1-latencyCurrentBeta)*float64(currentLatency))
	atomic.StoreInt64(&b.currentLatencyUs, newCurrent)
}

func (b *googleBreaker) updateLatency(latencyUs int64) {
	if latencyUs <= 0 || b.timeoutUs <= 0 {
		return
	}

	b.updateBaselineLatency(latencyUs)
	b.updateCurrentLatency(latencyUs)
}

type googlePromise struct {
	b     *googleBreaker
	start time.Duration
}

func (p googlePromise) Accept() {
	latencyUs := timex.Since(p.start).Microseconds()
	p.b.markSuccess(latencyUs)
}

func (p googlePromise) Reject() {
	p.b.markFailure()
}
