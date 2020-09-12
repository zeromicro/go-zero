package fx

import (
	"io/ioutil"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stringx"
)

func TestBuffer(t *testing.T) {
	const N = 5
	var count int32
	var wait sync.WaitGroup
	wait.Add(1)
	From(func(source chan<- interface{}) {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()

		for i := 0; i < 2*N; i++ {
			select {
			case source <- i:
				atomic.AddInt32(&count, 1)
			case <-ticker.C:
				wait.Done()
				return
			}
		}
	}).Buffer(N).ForAll(func(pipe <-chan interface{}) {
		wait.Wait()
		// why N+1, because take one more to wait for sending into the channel
		assert.Equal(t, int32(N+1), atomic.LoadInt32(&count))
	})
}

func TestBufferNegative(t *testing.T) {
	var result int
	Just(1, 2, 3, 4).Buffer(-1).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
		for item := range pipe {
			result += item.(int)
		}
		return result, nil
	})
	assert.Equal(t, 10, result)
}

func TestDone(t *testing.T) {
	var count int32
	Just(1, 2, 3).Walk(func(item interface{}, pipe chan<- interface{}) {
		time.Sleep(time.Millisecond * 100)
		atomic.AddInt32(&count, int32(item.(int)))
	}).Done()
	assert.Equal(t, int32(6), count)
}

func TestJust(t *testing.T) {
	var result int
	Just(1, 2, 3, 4).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
		for item := range pipe {
			result += item.(int)
		}
		return result, nil
	})
	assert.Equal(t, 10, result)
}

func TestDistinct(t *testing.T) {
	var result int
	Just(4, 1, 3, 2, 3, 4).Distinct(func(item interface{}) interface{} {
		return item
	}).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
		for item := range pipe {
			result += item.(int)
		}
		return result, nil
	})
	assert.Equal(t, 10, result)
}

func TestFilter(t *testing.T) {
	var result int
	Just(1, 2, 3, 4).Filter(func(item interface{}) bool {
		return item.(int)%2 == 0
	}).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
		for item := range pipe {
			result += item.(int)
		}
		return result, nil
	})
	assert.Equal(t, 6, result)
}

func TestForAll(t *testing.T) {
	var result int
	Just(1, 2, 3, 4).Filter(func(item interface{}) bool {
		return item.(int)%2 == 0
	}).ForAll(func(pipe <-chan interface{}) {
		for item := range pipe {
			result += item.(int)
		}
	})
	assert.Equal(t, 6, result)
}

func TestGroup(t *testing.T) {
	var groups [][]int
	Just(10, 11, 20, 21).Group(func(item interface{}) interface{} {
		v := item.(int)
		return v / 10
	}).ForEach(func(item interface{}) {
		v := item.([]interface{})
		var group []int
		for _, each := range v {
			group = append(group, each.(int))
		}
		groups = append(groups, group)
	})

	assert.Equal(t, 2, len(groups))
	for _, group := range groups {
		assert.Equal(t, 2, len(group))
		assert.True(t, group[0]/10 == group[1]/10)
	}
}

func TestHead(t *testing.T) {
	var result int
	Just(1, 2, 3, 4).Head(2).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
		for item := range pipe {
			result += item.(int)
		}
		return result, nil
	})
	assert.Equal(t, 3, result)
}

func TestHeadMore(t *testing.T) {
	var result int
	Just(1, 2, 3, 4).Head(6).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
		for item := range pipe {
			result += item.(int)
		}
		return result, nil
	})
	assert.Equal(t, 10, result)
}

func TestMap(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	tests := []struct {
		mapper MapFunc
		expect int
	}{
		{
			mapper: func(item interface{}) interface{} {
				v := item.(int)
				return v * v
			},
			expect: 30,
		},
		{
			mapper: func(item interface{}) interface{} {
				v := item.(int)
				if v%2 == 0 {
					return 0
				}
				return v * v
			},
			expect: 10,
		},
		{
			mapper: func(item interface{}) interface{} {
				v := item.(int)
				if v%2 == 0 {
					panic(v)
				}
				return v * v
			},
			expect: 10,
		},
	}

	// Map(...) works even WithWorkers(0)
	for i, test := range tests {
		t.Run(stringx.Rand(), func(t *testing.T) {
			var result int
			var workers int
			if i%2 == 0 {
				workers = 0
			} else {
				workers = runtime.NumCPU()
			}
			From(func(source chan<- interface{}) {
				for i := 1; i < 5; i++ {
					source <- i
				}
			}).Map(test.mapper, WithWorkers(workers)).Reduce(
				func(pipe <-chan interface{}) (interface{}, error) {
					for item := range pipe {
						result += item.(int)
					}
					return result, nil
				})

			assert.Equal(t, test.expect, result)
		})
	}
}

func TestMerge(t *testing.T) {
	Just(1, 2, 3, 4).Merge().ForEach(func(item interface{}) {
		assert.ElementsMatch(t, []interface{}{1, 2, 3, 4}, item.([]interface{}))
	})
}

func TestParallelJust(t *testing.T) {
	var count int32
	Just(1, 2, 3).Parallel(func(item interface{}) {
		time.Sleep(time.Millisecond * 100)
		atomic.AddInt32(&count, int32(item.(int)))
	}, UnlimitedWorkers())
	assert.Equal(t, int32(6), count)
}

func TestReverse(t *testing.T) {
	Just(1, 2, 3, 4).Reverse().Merge().ForEach(func(item interface{}) {
		assert.ElementsMatch(t, []interface{}{4, 3, 2, 1}, item.([]interface{}))
	})
}

func TestSort(t *testing.T) {
	var prev int
	Just(5, 3, 7, 1, 9, 6, 4, 8, 2).Sort(func(a, b interface{}) bool {
		return a.(int) < b.(int)
	}).ForEach(func(item interface{}) {
		next := item.(int)
		assert.True(t, prev < next)
		prev = next
	})
}

func TestTail(t *testing.T) {
	var result int
	Just(1, 2, 3, 4).Tail(2).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
		for item := range pipe {
			result += item.(int)
		}
		return result, nil
	})
	assert.Equal(t, 7, result)
}

func TestWalk(t *testing.T) {
	var result int
	Just(1, 2, 3, 4, 5).Walk(func(item interface{}, pipe chan<- interface{}) {
		if item.(int)%2 != 0 {
			pipe <- item
		}
	}, UnlimitedWorkers()).ForEach(func(item interface{}) {
		result += item.(int)
	})
	assert.Equal(t, 9, result)
}

func BenchmarkMapReduce(b *testing.B) {
	b.ReportAllocs()

	mapper := func(v interface{}) interface{} {
		return v.(int64) * v.(int64)
	}
	reducer := func(input <-chan interface{}) (interface{}, error) {
		var result int64
		for v := range input {
			result += v.(int64)
		}
		return result, nil
	}

	for i := 0; i < b.N; i++ {
		From(func(input chan<- interface{}) {
			for j := 0; j < 2; j++ {
				input <- int64(j)
			}
		}).Map(mapper).Reduce(reducer)
	}
}
