package iox

import (
	"bytes"
	"sync"
)

type BufferPool struct {
	capability int
	pool       *sync.Pool
}

func NewBufferPool(capability int) *BufferPool {
	return &BufferPool{
		capability: capability,
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

func (bp *BufferPool) Get() *bytes.Buffer {
	buf := bp.pool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func (bp *BufferPool) Put(buf *bytes.Buffer) {
	if buf.Cap() < bp.capability {
		bp.pool.Put(buf)
	}
}
