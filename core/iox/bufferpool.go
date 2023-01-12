package iox

import (
	"bytes"
	"sync"
)

// A BufferPool is a pool to buffer bytes.Buffer objects.
type BufferPool struct {
	capability int
	pool       *sync.Pool
}

// NewBufferPool returns a BufferPool.
func NewBufferPool(capability int) *BufferPool {
	return &BufferPool{
		capability: capability,
		pool: &sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
}

// Get returns a bytes.Buffer object from bp.
func (bp *BufferPool) Get() *bytes.Buffer {
	buf := bp.pool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// Put returns buf into bp.
func (bp *BufferPool) Put(buf *bytes.Buffer) {
	if buf.Cap() < bp.capability {
		bp.pool.Put(buf)
	}
}
