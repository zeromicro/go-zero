package syncx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtomicBool(t *testing.T) {
	val := ForAtomicBool(true)
	assert.True(t, val.True())
	val.Set(false)
	assert.False(t, val.True())
	val.Set(true)
	assert.True(t, val.True())
	val.Set(false)
	assert.False(t, val.True())
	ok := val.CompareAndSwap(false, true)
	assert.True(t, ok)
	assert.True(t, val.True())
	ok = val.CompareAndSwap(true, false)
	assert.True(t, ok)
	assert.False(t, val.True())
	ok = val.CompareAndSwap(true, false)
	assert.False(t, ok)
	assert.False(t, val.True())
}

func TestAtomicBoolFromVal(t *testing.T) {
	val := AtomicBoolFromVal(true)
	assert.True(t, val.Load())
	val.Store(false)
	assert.False(t, val.Load())
	val.Store(true)
	assert.True(t, val.Load())
	val.Store(false)
	assert.False(t, val.Load())
	ok := val.CompareAndSwap(false, true)
	assert.True(t, ok)
	assert.True(t, val.Load())
	ok = val.CompareAndSwap(true, false)
	assert.True(t, ok)
	assert.False(t, val.Load())
	ok = val.CompareAndSwap(true, false)
	assert.False(t, ok)
	assert.False(t, val.Load())
}
