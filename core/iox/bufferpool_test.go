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
