package collection

import (
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/timex"
)

const (
	testStep = time.Minute
	waitTime = time.Second
)

func TestNewTimingWheel(t *testing.T) {
	_, err := NewTimingWheel(0, 10, func(key, value any) {})
	assert.NotNil(t, err)
}

func TestTimingWheel_Drain(t *testing.T) {
	ticker := timex.NewFakeTicker()
	tw, _ := NewTimingWheelWithTicker(testStep, 10, func(k, v any) {
	}, ticker)
	tw.SetTimer("first", 3, testStep*4)
	tw.SetTimer("second", 5, testStep*7)
	tw.SetTimer("third", 7, testStep*7)
	var keys []string
	var vals []int
	var lock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(3)
	tw.Drain(func(key, value any) {
		lock.Lock()
		defer lock.Unlock()
		keys = append(keys, key.(string))
		vals = append(vals, value.(int))
		wg.Done()
	})
	wg.Wait()
	sort.Strings(keys)
	sort.Ints(vals)
	assert.Equal(t, 3, len(keys))
	assert.EqualValues(t, []string{"first", "second", "third"}, keys)
	assert.EqualValues(t, []int{3, 5, 7}, vals)
	var count int
	tw.Drain(func(key, value any) {
		count++
	})
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, 0, count)
	tw.Stop()
	assert.Equal(t, ErrClosed, tw.Drain(func(key, value any) {}))
}

func TestTimingWheel_SetTimerSoon(t *testing.T) {
	run := syncx.NewAtomicBool()
	ticker := timex.NewFakeTicker()
	tw, _ := NewTimingWheelWithTicker(testStep, 10, func(k, v any) {
		assert.True(t, run.CompareAndSwap(false, true))
		assert.Equal(t, "any", k)
		assert.Equal(t, 3, v.(int))
		ticker.Done()
	}, ticker)
	defer tw.Stop()
	tw.SetTimer("any", 3, testStep>>1)
	ticker.Tick()
	assert.Nil(t, ticker.Wait(waitTime))
	assert.True(t, run.True())
}

func TestTimingWheel_SetTimerTwice(t *testing.T) {
	run := syncx.NewAtomicBool()
	ticker := timex.NewFakeTicker()
	tw, _ := NewTimingWheelWithTicker(testStep, 10, func(k, v any) {
		assert.True(t, run.CompareAndSwap(false, true))
		assert.Equal(t, "any", k)
		assert.Equal(t, 5, v.(int))
		ticker.Done()
	}, ticker)
	defer tw.Stop()
	tw.SetTimer("any", 3, testStep*4)
	tw.SetTimer("any", 5, testStep*7)
	for i := 0; i < 8; i++ {
		ticker.Tick()
	}
	assert.Nil(t, ticker.Wait(waitTime))
	assert.True(t, run.True())
}

func TestTimingWheel_SetTimerWrongDelay(t *testing.T) {
	ticker := timex.NewFakeTicker()
	tw, _ := NewTimingWheelWithTicker(testStep, 10, func(k, v any) {}, ticker)
	defer tw.Stop()
	assert.NotPanics(t, func() {
		tw.SetTimer("any", 3, -testStep)
	})
}

func TestTimingWheel_SetTimerAfterClose(t *testing.T) {
	ticker := timex.NewFakeTicker()
	tw, _ := NewTimingWheelWithTicker(testStep, 10, func(k, v any) {}, ticker)
	tw.Stop()
	assert.Equal(t, ErrClosed, tw.SetTimer("any", 3, testStep))
}

func TestTimingWheel_MoveTimer(t *testing.T) {
	run := syncx.NewAtomicBool()
	ticker := timex.NewFakeTicker()
	tw, _ := NewTimingWheelWithTicker(testStep, 3, func(k, v any) {
		assert.True(t, run.CompareAndSwap(false, true))
		assert.Equal(t, "any", k)
		assert.Equal(t, 3, v.(int))
		ticker.Done()
	}, ticker)
	tw.SetTimer("any", 3, testStep*4)
	tw.MoveTimer("any", testStep*7)
	tw.MoveTimer("any", -testStep)
	tw.MoveTimer("none", testStep)
	for i := 0; i < 5; i++ {
		ticker.Tick()
	}
	assert.False(t, run.True())
	for i := 0; i < 3; i++ {
		ticker.Tick()
	}
	assert.Nil(t, ticker.Wait(waitTime))
	assert.True(t, run.True())
	tw.Stop()
	assert.Equal(t, ErrClosed, tw.MoveTimer("any", time.Millisecond))
}

func TestTimingWheel_MoveTimerSoon(t *testing.T) {
	run := syncx.NewAtomicBool()
	ticker := timex.NewFakeTicker()
	tw, _ := NewTimingWheelWithTicker(testStep, 3, func(k, v any) {
		assert.True(t, run.CompareAndSwap(false, true))
		assert.Equal(t, "any", k)
		assert.Equal(t, 3, v.(int))
		ticker.Done()
	}, ticker)
	defer tw.Stop()
	tw.SetTimer("any", 3, testStep*4)
	tw.MoveTimer("any", testStep>>1)
	assert.Nil(t, ticker.Wait(waitTime))
	assert.True(t, run.True())
}

func TestTimingWheel_MoveTimerEarlier(t *testing.T) {
	run := syncx.NewAtomicBool()
	ticker := timex.NewFakeTicker()
	tw, _ := NewTimingWheelWithTicker(testStep, 10, func(k, v any) {
		assert.True(t, run.CompareAndSwap(false, true))
		assert.Equal(t, "any", k)
		assert.Equal(t, 3, v.(int))
		ticker.Done()
	}, ticker)
	defer tw.Stop()
	tw.SetTimer("any", 3, testStep*4)
	tw.MoveTimer("any", testStep*2)
	for i := 0; i < 3; i++ {
		ticker.Tick()
	}
	assert.Nil(t, ticker.Wait(waitTime))
	assert.True(t, run.True())
}

func TestTimingWheel_RemoveTimer(t *testing.T) {
	ticker := timex.NewFakeTicker()
	tw, _ := NewTimingWheelWithTicker(testStep, 10, func(k, v any) {}, ticker)
	tw.SetTimer("any", 3, testStep)
	assert.NotPanics(t, func() {
		tw.RemoveTimer("any")
		tw.RemoveTimer("none")
		tw.RemoveTimer(nil)
	})
	for i := 0; i < 5; i++ {
		ticker.Tick()
	}
	tw.Stop()
	assert.Equal(t, ErrClosed, tw.RemoveTimer("any"))
}

func TestTimingWheel_SetTimer(t *testing.T) {
	tests := []struct {
		slots int
		setAt time.Duration
	}{
		{
			slots: 5,
			setAt: 5,
		},
		{
			slots: 5,
			setAt: 7,
		},
		{
			slots: 5,
			setAt: 10,
		},
		{
			slots: 5,
			setAt: 12,
		},
		{
			slots: 5,
			setAt: 7,
		},
		{
			slots: 5,
			setAt: 10,
		},
		{
			slots: 5,
			setAt: 12,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(stringx.RandId(), func(t *testing.T) {
			t.Parallel()

			var count int32
			ticker := timex.NewFakeTicker()
			tick := func() {
				atomic.AddInt32(&count, 1)
				ticker.Tick()
				time.Sleep(time.Millisecond)
			}
			var actual int32
			done := make(chan lang.PlaceholderType)
			tw, err := NewTimingWheelWithTicker(testStep, test.slots, func(key, value any) {
				assert.Equal(t, 1, key.(int))
				assert.Equal(t, 2, value.(int))
				actual = atomic.LoadInt32(&count)
				close(done)
			}, ticker)
			assert.Nil(t, err)
			defer tw.Stop()

			tw.SetTimer(1, 2, testStep*test.setAt)

			for {
				select {
				case <-done:
					assert.Equal(t, int32(test.setAt), actual)
					return
				default:
					tick()
				}
			}
		})
	}
}

func TestTimingWheel_SetAndMoveThenStart(t *testing.T) {
	tests := []struct {
		slots  int
		setAt  time.Duration
		moveAt time.Duration
	}{
		{
			slots:  5,
			setAt:  3,
			moveAt: 5,
		},
		{
			slots:  5,
			setAt:  3,
			moveAt: 7,
		},
		{
			slots:  5,
			setAt:  3,
			moveAt: 10,
		},
		{
			slots:  5,
			setAt:  3,
			moveAt: 12,
		},
		{
			slots:  5,
			setAt:  5,
			moveAt: 7,
		},
		{
			slots:  5,
			setAt:  5,
			moveAt: 10,
		},
		{
			slots:  5,
			setAt:  5,
			moveAt: 12,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(stringx.RandId(), func(t *testing.T) {
			t.Parallel()

			var count int32
			ticker := timex.NewFakeTicker()
			tick := func() {
				atomic.AddInt32(&count, 1)
				ticker.Tick()
				time.Sleep(time.Millisecond * 10)
			}
			var actual int32
			done := make(chan lang.PlaceholderType)
			tw, err := NewTimingWheelWithTicker(testStep, test.slots, func(key, value any) {
				actual = atomic.LoadInt32(&count)
				close(done)
			}, ticker)
			assert.Nil(t, err)
			defer tw.Stop()

			tw.SetTimer(1, 2, testStep*test.setAt)
			tw.MoveTimer(1, testStep*test.moveAt)

			for {
				select {
				case <-done:
					assert.Equal(t, int32(test.moveAt), actual)
					return
				default:
					tick()
				}
			}
		})
	}
}

func TestTimingWheel_SetAndMoveTwice(t *testing.T) {
	tests := []struct {
		slots       int
		setAt       time.Duration
		moveAt      time.Duration
		moveAgainAt time.Duration
	}{
		{
			slots:       5,
			setAt:       3,
			moveAt:      5,
			moveAgainAt: 10,
		},
		{
			slots:       5,
			setAt:       3,
			moveAt:      7,
			moveAgainAt: 12,
		},
		{
			slots:       5,
			setAt:       3,
			moveAt:      10,
			moveAgainAt: 15,
		},
		{
			slots:       5,
			setAt:       3,
			moveAt:      12,
			moveAgainAt: 17,
		},
		{
			slots:       5,
			setAt:       5,
			moveAt:      7,
			moveAgainAt: 12,
		},
		{
			slots:       5,
			setAt:       5,
			moveAt:      10,
			moveAgainAt: 17,
		},
		{
			slots:       5,
			setAt:       5,
			moveAt:      12,
			moveAgainAt: 17,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(stringx.RandId(), func(t *testing.T) {
			t.Parallel()

			var count int32
			ticker := timex.NewFakeTicker()
			tick := func() {
				atomic.AddInt32(&count, 1)
				ticker.Tick()
				time.Sleep(time.Millisecond * 10)
			}
			var actual int32
			done := make(chan lang.PlaceholderType)
			tw, err := NewTimingWheelWithTicker(testStep, test.slots, func(key, value any) {
				actual = atomic.LoadInt32(&count)
				close(done)
			}, ticker)
			assert.Nil(t, err)
			defer tw.Stop()

			tw.SetTimer(1, 2, testStep*test.setAt)
			tw.MoveTimer(1, testStep*test.moveAt)
			tw.MoveTimer(1, testStep*test.moveAgainAt)

			for {
				select {
				case <-done:
					assert.Equal(t, int32(test.moveAgainAt), actual)
					return
				default:
					tick()
				}
			}
		})
	}
}

func TestTimingWheel_ElapsedAndSet(t *testing.T) {
	tests := []struct {
		slots   int
		elapsed time.Duration
		setAt   time.Duration
	}{
		{
			slots:   5,
			elapsed: 3,
			setAt:   5,
		},
		{
			slots:   5,
			elapsed: 3,
			setAt:   7,
		},
		{
			slots:   5,
			elapsed: 3,
			setAt:   10,
		},
		{
			slots:   5,
			elapsed: 3,
			setAt:   12,
		},
		{
			slots:   5,
			elapsed: 5,
			setAt:   7,
		},
		{
			slots:   5,
			elapsed: 5,
			setAt:   10,
		},
		{
			slots:   5,
			elapsed: 5,
			setAt:   12,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(stringx.RandId(), func(t *testing.T) {
			t.Parallel()

			var count int32
			ticker := timex.NewFakeTicker()
			tick := func() {
				atomic.AddInt32(&count, 1)
				ticker.Tick()
				time.Sleep(time.Millisecond * 10)
			}
			var actual int32
			done := make(chan lang.PlaceholderType)
			tw, err := NewTimingWheelWithTicker(testStep, test.slots, func(key, value any) {
				actual = atomic.LoadInt32(&count)
				close(done)
			}, ticker)
			assert.Nil(t, err)
			defer tw.Stop()

			for i := 0; i < int(test.elapsed); i++ {
				tick()
			}

			tw.SetTimer(1, 2, testStep*test.setAt)

			for {
				select {
				case <-done:
					assert.Equal(t, int32(test.elapsed+test.setAt), actual)
					return
				default:
					tick()
				}
			}
		})
	}
}

func TestTimingWheel_ElapsedAndSetThenMove(t *testing.T) {
	tests := []struct {
		slots   int
		elapsed time.Duration
		setAt   time.Duration
		moveAt  time.Duration
	}{
		{
			slots:   5,
			elapsed: 3,
			setAt:   5,
			moveAt:  10,
		},
		{
			slots:   5,
			elapsed: 3,
			setAt:   7,
			moveAt:  12,
		},
		{
			slots:   5,
			elapsed: 3,
			setAt:   10,
			moveAt:  15,
		},
		{
			slots:   5,
			elapsed: 3,
			setAt:   12,
			moveAt:  16,
		},
		{
			slots:   5,
			elapsed: 5,
			setAt:   7,
			moveAt:  12,
		},
		{
			slots:   5,
			elapsed: 5,
			setAt:   10,
			moveAt:  15,
		},
		{
			slots:   5,
			elapsed: 5,
			setAt:   12,
			moveAt:  17,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(stringx.RandId(), func(t *testing.T) {
			t.Parallel()

			var count int32
			ticker := timex.NewFakeTicker()
			tick := func() {
				atomic.AddInt32(&count, 1)
				ticker.Tick()
				time.Sleep(time.Millisecond * 10)
			}
			var actual int32
			done := make(chan lang.PlaceholderType)
			tw, err := NewTimingWheelWithTicker(testStep, test.slots, func(key, value any) {
				actual = atomic.LoadInt32(&count)
				close(done)
			}, ticker)
			assert.Nil(t, err)
			defer tw.Stop()

			for i := 0; i < int(test.elapsed); i++ {
				tick()
			}

			tw.SetTimer(1, 2, testStep*test.setAt)
			tw.MoveTimer(1, testStep*test.moveAt)

			for {
				select {
				case <-done:
					assert.Equal(t, int32(test.elapsed+test.moveAt), actual)
					return
				default:
					tick()
				}
			}
		})
	}
}

func TestMoveAndRemoveTask(t *testing.T) {
	ticker := timex.NewFakeTicker()
	tick := func(v int) {
		for i := 0; i < v; i++ {
			ticker.Tick()
		}
	}
	var keys []int
	tw, _ := NewTimingWheelWithTicker(testStep, 10, func(k, v any) {
		assert.Equal(t, "any", k)
		assert.Equal(t, 3, v.(int))
		keys = append(keys, v.(int))
		ticker.Done()
	}, ticker)
	defer tw.Stop()
	tw.SetTimer("any", 3, testStep*8)
	tick(6)
	tw.MoveTimer("any", testStep*7)
	tick(3)
	tw.RemoveTimer("any")
	tick(30)
	time.Sleep(time.Millisecond)
	assert.Equal(t, 0, len(keys))
}

func BenchmarkTimingWheel(b *testing.B) {
	b.ReportAllocs()

	tw, _ := NewTimingWheel(time.Second, 100, func(k, v any) {})
	for i := 0; i < b.N; i++ {
		tw.SetTimer(i, i, time.Second)
		tw.SetTimer(b.N+i, b.N+i, time.Second)
		tw.MoveTimer(i, time.Second*time.Duration(i))
		tw.RemoveTimer(i)
	}
}
