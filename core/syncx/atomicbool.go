package syncx

import "sync/atomic"

// An AtomicBool is an atomic implementation for boolean values.
type AtomicBool struct {
	b atomic.Bool
}

// NewAtomicBool returns an AtomicBool.
// Deprecated: use atomic.Bool instead.
func NewAtomicBool() *AtomicBool {
	return &AtomicBool{
		b: atomic.Bool{},
	}
}

// ForAtomicBool returns an AtomicBool with given val.
func ForAtomicBool(val bool) *AtomicBool {
	b := NewAtomicBool()
	b.Set(val)
	return b
}

// CompareAndSwap compares current value with given old, if equals, set to given val.
func (b *AtomicBool) CompareAndSwap(old, val bool) bool {
	return b.b.CompareAndSwap(old, val)
}

// Set sets the value to v.
func (b *AtomicBool) Set(v bool) {
	b.b.Store(v)
}

// True returns true if current value is true.
func (b *AtomicBool) True() bool {
	return b.b.Load()
}
