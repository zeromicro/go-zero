package bytesx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
)

func TestMallocSize(t *testing.T) {
	n := 100
	waitGroup := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		waitGroup.Add(1)
		go func(j int) {
			defer waitGroup.Done()
			k := j % 10
			t.Run(fmt.Sprintf("%d bytes", k), func(t *testing.T) {
				b := MallocSize(k)
				assert.EqualValues(t, len(b), k)
				Free(b)
			})

		}(i)
	}
	waitGroup.Wait()
}

func TestMalloc(t *testing.T) {
	b := Malloc(10, 20)
	assert.EqualValues(t, 10, len(b))
	assert.EqualValues(t, 20, cap(b))
	Free(b)
}

func BenchmarkMakeBytes(b *testing.B) {
	b.ReportAllocs()
	rand.Seed(2022)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			waitGroup := sync.WaitGroup{}
			for j := 0; j < 20; j++ {
				b.StopTimer()
				waitGroup.Add(1)
				k := 2 << j
				b.StartTimer()
				go func() {
					a := make([]byte, k)
					_ = a
					b.StopTimer()
					waitGroup.Done()
					b.StartTimer()
				}()
			}
			waitGroup.Wait()

		}
	})

}

func BenchmarkMalloc(b *testing.B) {
	rand.Seed(2022)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			waitGroup := sync.WaitGroup{}
			for j := 0; j < 20; j++ {
				b.StopTimer()
				k := 2 << j
				waitGroup.Add(1)
				b.StartTimer()
				go func() {
					a := MallocSize(k)
					b.StopTimer()
					Free(a)
					waitGroup.Done()
					b.StartTimer()
				}()
			}
			waitGroup.Wait()

		}
	})

}
