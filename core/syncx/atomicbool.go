package syncx

import (
	"sync/atomic"
)

// An AtomicBool is an atomic implementation for boolean values.
type AtomicBool struct {
	_ noCopy //no copy
	v uint32
}

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

// NewAtomicBool returns an AtomicBool.
// Deprecated: use atomic.Bool instead.
func NewAtomicBool() AtomicBool {
	return AtomicBool{
		v: 0,
	}
}

// ForAtomicBool returns an AtomicBool with given val.
// Deprecated: use atomic.Bool instead.
func ForAtomicBool(val bool) AtomicBool {
	return AtomicBool{
		v: b32(val),
	}
}

// CompareAndSwap compares current value with given old, if equals, set to given val.
func (b *AtomicBool) CompareAndSwap(old, val bool) bool {
	return atomic.CompareAndSwapUint32(&b.v, b32(old), b32(val))
}

// Set sets the value to v.
func (b *AtomicBool) Set(v bool) {
	atomic.StoreUint32(&b.v, b32(v))
}

// True returns true if current value is true.
func (b *AtomicBool) True() bool {
	return atomic.LoadUint32(&b.v) != 0
}

// b32 returns a uint32 0 or 1 representing b.
func b32(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}
