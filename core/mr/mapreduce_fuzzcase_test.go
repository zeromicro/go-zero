//go:build fuzz

package mr

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/threading"
	"gopkg.in/cheggaaa/pb.v1"
)

// If Fuzz stuck, we don't know why, because it only returns hung or unexpected,
// so we need to simulate the fuzz test in test mode.
func TestMapReduceRandom(t *testing.T) {
	rand.NewSource(time.Now().UnixNano())

	const (
		times  = 10000
		nRange = 500
		mega   = 1024 * 1024
	)

	bar := pb.New(times).Start()
	runner := threading.NewTaskRunner(runtime.NumCPU())
	var wg sync.WaitGroup
	wg.Add(times)
	for i := 0; i < times; i++ {
		runner.Schedule(func() {
			start := time.Now()
			defer func() {
				if time.Since(start) > time.Minute {
					t.Fatal("timeout")
				}
				wg.Done()
			}()

			t.Run(strconv.Itoa(i), func(t *testing.T) {
				n := rand.Int63n(nRange)%nRange + nRange
				workers := rand.Int()%50 + runtime.NumCPU()/2
				genPanic := rand.Intn(100) == 0
				mapperPanic := rand.Intn(100) == 0
				reducerPanic := rand.Intn(100) == 0
				genIdx := rand.Int63n(n)
				mapperIdx := rand.Int63n(n)
				reducerIdx := rand.Int63n(n)
				squareSum := (n - 1) * n * (2*n - 1) / 6

				fn := func() (int64, error) {
					return MapReduce(func(source chan<- int64) {
						for i := int64(0); i < n; i++ {
							source <- i
							if genPanic && i == genIdx {
								panic("foo")
							}
						}
					}, func(v int64, writer Writer[int64], cancel func(error)) {
						if mapperPanic && v == mapperIdx {
							panic("bar")
						}
						writer.Write(v * v)
					}, func(pipe <-chan int64, writer Writer[int64], cancel func(error)) {
						var idx int64
						var total int64
						for v := range pipe {
							if reducerPanic && idx == reducerIdx {
								panic("baz")
							}
							total += v
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
					assert.Equal(t, squareSum, val)
				}
				bar.Increment()
			})
		})
	}

	wg.Wait()
	bar.Finish()
}
