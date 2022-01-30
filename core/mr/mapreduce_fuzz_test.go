//go:build go1.18
// +build go1.18

package mr

import (
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func FuzzMapReduce(f *testing.F) {
	rand.Seed(time.Now().UnixNano())

	f.Add(uint(10), uint(runtime.NumCPU()))
	f.Fuzz(func(t *testing.T, num uint, workers uint) {
		n := int64(num)%5000 + 5000
		genPanic := rand.Intn(100) == 0
		mapperPanic := rand.Intn(100) == 0
		reducerPanic := rand.Intn(100) == 0
		genIdx := rand.Int63n(n)
		mapperIdx := rand.Int63n(n)
		reducerIdx := rand.Int63n(n)
		squareSum := (n - 1) * n * (2*n - 1) / 6

		fn := func() (interface{}, error) {
			defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

			return MapReduce(func(source chan<- interface{}) {
				for i := int64(0); i < n; i++ {
					source <- i
					if genPanic && i == genIdx {
						panic("foo")
					}
				}
			}, func(item interface{}, writer Writer, cancel func(error)) {
				v := item.(int64)
				if mapperPanic && v == mapperIdx {
					panic("bar")
				}
				writer.Write(v * v)
			}, func(pipe <-chan interface{}, writer Writer, cancel func(error)) {
				var idx int64
				var total int64
				for v := range pipe {
					if reducerPanic && idx == reducerIdx {
						panic("baz")
					}
					total += v.(int64)
					idx++
				}
				writer.Write(total)
			}, WithWorkers(int(workers)%50+runtime.NumCPU()/2))
		}

		if genPanic || mapperPanic || reducerPanic {
			var buf strings.Builder
			buf.WriteString(fmt.Sprintf("n: %d", n))
			buf.WriteString(fmt.Sprintf(", genPanic: %t", genPanic))
			buf.WriteString(fmt.Sprintf(", mapperPanic: %t", mapperPanic))
			buf.WriteString(fmt.Sprintf(", reducerPanic: %t", reducerPanic))
			buf.WriteString(fmt.Sprintf(", genIdx: %d", genIdx))
			buf.WriteString(fmt.Sprintf(", mapperIdx: %d", mapperIdx))
			buf.WriteString(fmt.Sprintf(", reducerIdx: %d", reducerIdx))
			assert.Panicsf(t, func() { fn() }, buf.String())
		} else {
			val, err := fn()
			assert.Nil(t, err)
			assert.Equal(t, squareSum, val.(int64))
		}
	})
}
