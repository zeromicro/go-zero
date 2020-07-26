package syncx

import (
	"math"
	"sync/atomic"
)

type AtomicFloat64 uint64

func NewAtomicFloat64() *AtomicFloat64 {
	return new(AtomicFloat64)
}

func ForAtomicFloat64(val float64) *AtomicFloat64 {
	f := NewAtomicFloat64()
	f.Set(val)
	return f
}

func (f *AtomicFloat64) Add(val float64) float64 {
	for {
		old := f.Load()
		nv := old + val
		if f.CompareAndSwap(old, nv) {
			return nv
		}
	}
}

func (f *AtomicFloat64) CompareAndSwap(old, val float64) bool {
	return atomic.CompareAndSwapUint64((*uint64)(f), math.Float64bits(old), math.Float64bits(val))
}

func (f *AtomicFloat64) Load() float64 {
	return math.Float64frombits(atomic.LoadUint64((*uint64)(f)))
}

func (f *AtomicFloat64) Set(val float64) {
	atomic.StoreUint64((*uint64)(f), math.Float64bits(val))
}
