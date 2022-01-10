package fx

import (
	"io/ioutil"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stringx"
	"go.uber.org/goleak"
)

func TestBuffer(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
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
	})
}

func TestBufferNegative(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result int
		Just(1, 2, 3, 4).Buffer(-1).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
			for item := range pipe {
				result += item.(int)
			}
			return result, nil
		})
		assert.Equal(t, 10, result)
	})
}

func TestCount(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		tests := []struct {
			name     string
			elements []interface{}
		}{
			{
				name: "no elements with nil",
			},
			{
				name:     "no elements",
				elements: []interface{}{},
			},
			{
				name:     "1 element",
				elements: []interface{}{1},
			},
			{
				name:     "multiple elements",
				elements: []interface{}{1, 2, 3},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				val := Just(test.elements...).Count()
				assert.Equal(t, len(test.elements), val)
			})
		}
	})
}

func TestDone(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var count int32
		Just(1, 2, 3).Walk(func(item interface{}, pipe chan<- interface{}) {
			time.Sleep(time.Millisecond * 100)
			atomic.AddInt32(&count, int32(item.(int)))
		}).Done()
		assert.Equal(t, int32(6), count)
	})
}

func TestJust(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result int
		Just(1, 2, 3, 4).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
			for item := range pipe {
				result += item.(int)
			}
			return result, nil
		})
		assert.Equal(t, 10, result)
	})
}

func TestDistinct(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
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
	})
}

func TestFilter(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
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
	})
}

func TestFirst(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		assert.Nil(t, Just().First())
		assert.Equal(t, "foo", Just("foo").First())
		assert.Equal(t, "foo", Just("foo", "bar").First())
	})
}

func TestForAll(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result int
		Just(1, 2, 3, 4).Filter(func(item interface{}) bool {
			return item.(int)%2 == 0
		}).ForAll(func(pipe <-chan interface{}) {
			for item := range pipe {
				result += item.(int)
			}
		})
		assert.Equal(t, 6, result)
	})
}

func TestGroup(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
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
	})
}

func TestHead(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result int
		Just(1, 2, 3, 4).Head(2).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
			for item := range pipe {
				result += item.(int)
			}
			return result, nil
		})
		assert.Equal(t, 3, result)
	})
}

func TestHeadZero(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		assert.Panics(t, func() {
			Just(1, 2, 3, 4).Head(0).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
				return nil, nil
			})
		})
	})
}

func TestHeadMore(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result int
		Just(1, 2, 3, 4).Head(6).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
			for item := range pipe {
				result += item.(int)
			}
			return result, nil
		})
		assert.Equal(t, 10, result)
	})
}

func TestLast(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		goroutines := runtime.NumGoroutine()
		assert.Nil(t, Just().Last())
		assert.Equal(t, "foo", Just("foo").Last())
		assert.Equal(t, "bar", Just("foo", "bar").Last())
		// let scheduler schedule first
		runtime.Gosched()
		assert.Equal(t, goroutines, runtime.NumGoroutine())
	})
}

func TestMap(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
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
	})
}

func TestMerge(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		Just(1, 2, 3, 4).Merge().ForEach(func(item interface{}) {
			assert.ElementsMatch(t, []interface{}{1, 2, 3, 4}, item.([]interface{}))
		})
	})
}

func TestParallelJust(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var count int32
		Just(1, 2, 3).Parallel(func(item interface{}) {
			time.Sleep(time.Millisecond * 100)
			atomic.AddInt32(&count, int32(item.(int)))
		}, UnlimitedWorkers())
		assert.Equal(t, int32(6), count)
	})
}

func TestReverse(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		Just(1, 2, 3, 4).Reverse().Merge().ForEach(func(item interface{}) {
			assert.ElementsMatch(t, []interface{}{4, 3, 2, 1}, item.([]interface{}))
		})
	})
}

func TestSort(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var prev int
		Just(5, 3, 7, 1, 9, 6, 4, 8, 2).Sort(func(a, b interface{}) bool {
			return a.(int) < b.(int)
		}).ForEach(func(item interface{}) {
			next := item.(int)
			assert.True(t, prev < next)
			prev = next
		})
	})
}

func TestSplit(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		assert.Panics(t, func() {
			Just(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Split(0).Done()
		})
		var chunks [][]interface{}
		Just(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Split(4).ForEach(func(item interface{}) {
			chunk := item.([]interface{})
			chunks = append(chunks, chunk)
		})
		assert.EqualValues(t, [][]interface{}{
			{1, 2, 3, 4},
			{5, 6, 7, 8},
			{9, 10},
		}, chunks)
	})
}

func TestTail(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result int
		Just(1, 2, 3, 4).Tail(2).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
			for item := range pipe {
				result += item.(int)
			}
			return result, nil
		})
		assert.Equal(t, 7, result)
	})
}

func TestTailZero(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		assert.Panics(t, func() {
			Just(1, 2, 3, 4).Tail(0).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
				return nil, nil
			})
		})
	})
}

func TestWalk(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		var result int
		Just(1, 2, 3, 4, 5).Walk(func(item interface{}, pipe chan<- interface{}) {
			if item.(int)%2 != 0 {
				pipe <- item
			}
		}, UnlimitedWorkers()).ForEach(func(item interface{}) {
			result += item.(int)
		})
		assert.Equal(t, 9, result)
	})
}

func TestStream_AnyMach(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		assetEqual(t, false, Just(1, 2, 3).AnyMach(func(item interface{}) bool {
			return item.(int) == 4
		}))
		assetEqual(t, false, Just(1, 2, 3).AnyMach(func(item interface{}) bool {
			return item.(int) == 0
		}))
		assetEqual(t, true, Just(1, 2, 3).AnyMach(func(item interface{}) bool {
			return item.(int) == 2
		}))
		assetEqual(t, true, Just(1, 2, 3).AnyMach(func(item interface{}) bool {
			return item.(int) == 2
		}))
	})
}

func TestStream_AllMach(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		assetEqual(
			t, true, Just(1, 2, 3).AllMach(func(item interface{}) bool {
				return true
			}),
		)
		assetEqual(
			t, false, Just(1, 2, 3).AllMach(func(item interface{}) bool {
				return false
			}),
		)
		assetEqual(
			t, false, Just(1, 2, 3).AllMach(func(item interface{}) bool {
				return item.(int) == 1
			}),
		)
	})
}

func TestStream_NoneMatch(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		assetEqual(
			t, true, Just(1, 2, 3).NoneMatch(func(item interface{}) bool {
				return false
			}),
		)
		assetEqual(
			t, false, Just(1, 2, 3).NoneMatch(func(item interface{}) bool {
				return true
			}),
		)
		assetEqual(
			t, true, Just(1, 2, 3).NoneMatch(func(item interface{}) bool {
				return item.(int) == 4
			}),
		)
	})
}

func TestConcat(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		a1 := []interface{}{1, 2, 3}
		a2 := []interface{}{4, 5, 6}
		s1 := Just(a1...)
		s2 := Just(a2...)
		stream := Concat(s1, s2)
		var items []interface{}
		for item := range stream.source {
			items = append(items, item)
		}
		sort.Slice(items, func(i, j int) bool {
			return items[i].(int) < items[j].(int)
		})
		ints := make([]interface{}, 0)
		ints = append(ints, a1...)
		ints = append(ints, a2...)
		assetEqual(t, ints, items)
	})
}

func TestStream_Skip(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		assetEqual(t, 3, Just(1, 2, 3, 4).Skip(1).Count())
		assetEqual(t, 1, Just(1, 2, 3, 4).Skip(3).Count())
		assetEqual(t, 4, Just(1, 2, 3, 4).Skip(0).Count())
		equal(t, Just(1, 2, 3, 4).Skip(3), []interface{}{4})
		assert.Panics(t, func() {
			Just(1, 2, 3, 4).Skip(-1)
		})
	})
}

func TestStream_Concat(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		stream := Just(1).Concat(Just(2), Just(3))
		var items []interface{}
		for item := range stream.source {
			items = append(items, item)
		}
		sort.Slice(items, func(i, j int) bool {
			return items[i].(int) < items[j].(int)
		})
		assetEqual(t, []interface{}{1, 2, 3}, items)

		just := Just(1)
		equal(t, just.Concat(just), []interface{}{1})
	})
}

func BenchmarkParallelMapReduce(b *testing.B) {
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
	b.ResetTimer()
	From(func(input chan<- interface{}) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				input <- int64(rand.Int())
			}
		})
	}).Map(mapper).Reduce(reducer)
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
	b.ResetTimer()
	From(func(input chan<- interface{}) {
		for i := 0; i < b.N; i++ {
			input <- int64(rand.Int())
		}
	}).Map(mapper).Reduce(reducer)
}

func assetEqual(t *testing.T, except, data interface{}) {
	if !reflect.DeepEqual(except, data) {
		t.Errorf(" %v, want %v", data, except)
	}
}

func equal(t *testing.T, stream Stream, data []interface{}) {
	items := make([]interface{}, 0)
	for item := range stream.source {
		items = append(items, item)
	}
	if !reflect.DeepEqual(items, data) {
		t.Errorf(" %v, want %v", items, data)
	}
}

func runCheckedTest(t *testing.T, fn func(t *testing.T)) {
	defer goleak.VerifyNone(t)
	fn(t)
}
