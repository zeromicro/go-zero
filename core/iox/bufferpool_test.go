package iox

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBufferPool(t *testing.T) {
	capacity := 1024
	pool := NewBufferPool(capacity)
	pool.Put(bytes.NewBuffer(make([]byte, 0, 2*capacity)))
	assert.True(t, pool.Get().Cap() <= capacity)
}

func TestBufferPool_Put(t *testing.T) {
	t.Run("with nil buf", func(t *testing.T) {
		pool := NewBufferPool(1024)
		pool.Put(nil)
		val := pool.Get()
		assert.IsType(t, new(bytes.Buffer), val)
	})

	t.Run("with less-cap buf", func(t *testing.T) {
		pool := NewBufferPool(1024)
		pool.Put(bytes.NewBuffer(make([]byte, 0, 512)))
		val := pool.Get()
		assert.IsType(t, new(bytes.Buffer), val)
	})

	t.Run("with more-cap buf", func(t *testing.T) {
		pool := NewBufferPool(1024)
		pool.Put(bytes.NewBuffer(make([]byte, 0, 1024<<1)))
		val := pool.Get()
		assert.IsType(t, new(bytes.Buffer), val)
	})
}
