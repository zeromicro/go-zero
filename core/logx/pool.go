package logx

import (
	"strings"
	"sync"
)

var (
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
)

func getBuffer() *strings.Builder {
	return bufferPool.Get().(*strings.Builder)
}

func putBuffer(b *strings.Builder) {
	b.Reset()
	bufferPool.Put(b)
}
