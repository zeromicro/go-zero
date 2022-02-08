package bytesx

import (
	"github.com/zeromicro/go-zero/core/syncx"
	"strconv"
	"sync"
)

var (
	bytesPoolMap = &sync.Map{}
	singleFlight = syncx.NewSingleFlight()
)

// Malloc returns a bytes of slice.
func Malloc(size, capacity int) []byte {
	c := size
	if capacity > size {
		c = capacity
	}

	pool := getOrCreatePool(c)
	s := pool.Get()
	data := *(s.(*[]byte))

	return data[:size]
}

// MallocSize returns a bytes of slice.
func MallocSize(size int) []byte {
	return Malloc(size, size)
}

// Free recovers a bytes of slice.
func Free(buf []byte) {
	c := cap(buf)
	pool := getOrCreatePool(c)
	pool.Put(&buf)
}

func getOrCreatePool(c int) *sync.Pool {
	pool, _ := singleFlight.Do(strconv.Itoa(c), func() (interface{}, error) {
		if pool, ok := bytesPoolMap.Load(c); !ok {
			p := &sync.Pool{New: func() interface{} {
				s := make([]byte, 0, c)
				return &s
			}}
			bytesPoolMap.Store(c, p)
			return p, nil
		} else {
			return pool, nil
		}
	})
	return pool.(*sync.Pool)
}
