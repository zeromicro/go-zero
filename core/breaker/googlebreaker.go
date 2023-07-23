package breaker

import (
	"math"
	"time"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/mathx"
)

const (
	// 250ms for bucket duration
	window     = time.Second * 10
	buckets    = 40
	k          = 1.5
	protection = 5
)

// googleBreaker is a netflixBreaker pattern from google.
// see Client-Side Throttling section in https://landing.google.com/sre/sre-book/chapters/handling-overload/
type googleBreaker struct {
	k     float64
	stat  *collection.RollingWindow
	proba *mathx.Proba
}

func newGoogleBreaker() *googleBreaker {
	bucketDuration := time.Duration(int64(window) / int64(buckets))
	st := collection.NewRollingWindow(buckets, bucketDuration)
	return &googleBreaker{
		stat:  st,
		k:     k,
		proba: mathx.NewProba(),
	}
}

func (b *googleBreaker) accept() error {
	accepts, total := b.history()
	weightedAccepts := b.k * float64(accepts)
	// https://landing.google.com/sre/sre-book/chapters/handling-overload/#eq2101
	dropRatio := math.Max(0, (float64(total-protection)-weightedAccepts)/float64(total+1))
	if dropRatio <= 0 {
		return nil
	}

	if b.proba.TrueOnProba(dropRatio) {
		return ErrServiceUnavailable
	}

	return nil
}

func (b *googleBreaker) allow() (internalPromise, error) {
	if err := b.accept(); err != nil {
		return nil, err
	}

	return googlePromise{
		b: b,
	}, nil
}

func (b *googleBreaker) doReq(req func() error, fallback func(err error) error, acceptable Acceptable) error {
	if err := b.accept(); err != nil {
		if fallback != nil {
			return fallback(err)
		}

		return err
	}

	defer func() {
		if e := recover(); e != nil {
			b.markFailure()
			panic(e)
		}
	}()

	err := req()
	if acceptable(err) {
		b.markSuccess()
	} else {
		b.markFailure()
	}

	return err
}

func (b *googleBreaker) markSuccess() {
	b.stat.Add(1)
}

func (b *googleBreaker) markFailure() {
	b.stat.Add(0)
}

func (b *googleBreaker) history() (accepts, total int64) {
	b.stat.Reduce(func(b *collection.Bucket) {
		accepts += int64(b.Sum)
		total += b.Count
	})

	return
}

type googlePromise struct {
	b *googleBreaker
}

func (p googlePromise) Accept() {
	p.b.markSuccess()
}

func (p googlePromise) Reject() {
	p.b.markFailure()
}
