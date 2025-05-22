package syncx

import (
	"sync/atomic"
)

// An AtomicBool is an atomic implementation for boolean values.
type AtomicBool atomic.Bool

// NewAtomicBool returns an AtomicBool.
// Deprecated: use atomic.Bool instead.
func NewAtomicBool() (d AtomicBool) {
	return
}

// ForAtomicBool returns an AtomicBool with given val.
// Deprecated: use atomic.Bool instead.
func ForAtomicBool(val bool) (b AtomicBool) {
	b.Set(val)
	return
}

// CompareAndSwap compares current value with given old, if equals, set to given val.
func (b *AtomicBool) CompareAndSwap(old, val bool) bool {
	return (*atomic.Bool)(b).CompareAndSwap(old, val)
}

// Set sets the value to v.
func (b *AtomicBool) Set(v bool) {
	(*atomic.Bool)(b).Store(v)
}

// True returns true if current value is true.
func (b *AtomicBool) True() bool {
	return (*atomic.Bool)(b).Load()
}

// AtomicBoolFromVal returns an atomic.Bool with given val.
func AtomicBoolFromVal(val bool) (b atomic.Bool) {
	b.Store(val)
	return
}
