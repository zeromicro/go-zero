package syncx

import "sync/atomic"

// An AtomicBool is an atomic implementation for boolean values.
type AtomicBool uint32

// NewAtomicBool returns an AtomicBool.
func NewAtomicBool() *AtomicBool {
	return new(AtomicBool)
}

// ForAtomicBool returns an AtomicBool with given val.
func ForAtomicBool(val bool) *AtomicBool {
	b := NewAtomicBool()
	b.Set(val)
	return b
}

// CompareAndSwap compares current value with given old, if equals, set to given val.
func (b *AtomicBool) CompareAndSwap(old, val bool) bool {
	var ov, nv uint32

	if old {
		ov = 1
	}
	if val {
		nv = 1
	}

	return atomic.CompareAndSwapUint32((*uint32)(b), ov, nv)
}

// Set sets the value to v.
func (b *AtomicBool) Set(v bool) {
	if v {
		atomic.StoreUint32((*uint32)(b), 1)
	} else {
		atomic.StoreUint32((*uint32)(b), 0)
	}
}

// True returns true if current value is true.
func (b *AtomicBool) True() bool {
	return atomic.LoadUint32((*uint32)(b)) == 1
}
