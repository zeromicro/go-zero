package mathx

import (
	"math/rand"
	"sync"
	"time"
)

// An Unstable is used to generate random value around the mean value base on given deviation.
type Unstable struct {
	deviation float64
	p         *sync.Pool
}

// NewUnstable returns an Unstable.
func NewUnstable(deviation float64) Unstable {
	if deviation < 0 {
		deviation = 0
	}
	if deviation > 1 {
		deviation = 1
	}
	return Unstable{
		deviation: deviation,
		p: &sync.Pool{
			New: func() any {
				return rand.New(rand.NewSource(time.Now().UnixNano()))
			},
		},
	}
}

// AroundDuration returns a random duration with given base and deviation.
func (u Unstable) AroundDuration(base time.Duration) time.Duration {
	r := u.p.Get().(*rand.Rand)
	val := time.Duration((1 + u.deviation - 2*u.deviation*r.Float64()) * float64(base))
	u.p.Put(r)
	return val
}

// AroundInt returns a random int64 with given base and deviation.
func (u Unstable) AroundInt(base int64) int64 {
	r := u.p.Get().(*rand.Rand)
	val := int64((1 + u.deviation - 2*u.deviation*r.Float64()) * float64(base))
	u.p.Put(r)
	return val
}
