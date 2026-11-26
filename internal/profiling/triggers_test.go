package profiling

import (
	"math"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsOverloaded(t *testing.T) {
	tests := []struct {
		name           string
		cpuThreshold   int64
		heapThreshold  uint64
		goroutineLimit int
		heap           uint64
		goroutines     int
		want           bool
		heapReads      int
		goroutineReads int
	}{
		{name: "CPU zero retains unconditional profiling", cpuThreshold: 0,
			heapThreshold: 100, goroutineLimit: 10, want: true},
		{name: "new triggers disabled by default", cpuThreshold: math.MaxInt64,
			heap: math.MaxUint64, goroutines: math.MaxInt},
		{name: "heap below threshold", cpuThreshold: math.MaxInt64,
			heapThreshold: 100, heap: 99, heapReads: 1},
		{name: "heap at threshold", cpuThreshold: math.MaxInt64,
			heapThreshold: 100, heap: 100, want: true, heapReads: 1},
		{name: "heap above threshold", cpuThreshold: math.MaxInt64,
			heapThreshold: 100, heap: 101, want: true, heapReads: 1},
		{name: "goroutines below threshold", cpuThreshold: math.MaxInt64,
			goroutineLimit: 10, goroutines: 9, goroutineReads: 1},
		{name: "goroutines at threshold", cpuThreshold: math.MaxInt64,
			goroutineLimit: 10, goroutines: 10, want: true, goroutineReads: 1},
		{name: "goroutines above threshold", cpuThreshold: math.MaxInt64,
			goroutineLimit: 10, goroutines: 11, want: true, goroutineReads: 1},
		{name: "heap alone triggers with both enabled", cpuThreshold: math.MaxInt64,
			heapThreshold: 100, goroutineLimit: 10, heap: 100, goroutines: 9,
			want: true, heapReads: 1},
		{name: "goroutines alone trigger with both enabled", cpuThreshold: math.MaxInt64,
			heapThreshold: 100, goroutineLimit: 10, heap: 99, goroutines: 10,
			want: true, heapReads: 1, goroutineReads: 1},
		{name: "both below threshold", cpuThreshold: math.MaxInt64,
			heapThreshold: 100, goroutineLimit: 10, heap: 99, goroutines: 9,
			heapReads: 1, goroutineReads: 1},
		{name: "both at threshold need one trigger", cpuThreshold: math.MaxInt64,
			heapThreshold: 100, goroutineLimit: 10, heap: 100, goroutines: 10,
			want: true, heapReads: 1},
		{name: "disabled heap does not trigger", cpuThreshold: math.MaxInt64,
			goroutineLimit: 10, heap: math.MaxUint64, goroutines: 9, goroutineReads: 1},
		{name: "disabled goroutines do not trigger", cpuThreshold: math.MaxInt64,
			heapThreshold: 100, heap: 99, goroutines: math.MaxInt, heapReads: 1},
		{name: "large heap threshold does not overflow", cpuThreshold: math.MaxInt64,
			heapThreshold: 1 << 63, heap: 1<<63 - 1, heapReads: 1},
		{name: "large heap value triggers", cpuThreshold: math.MaxInt64,
			heapThreshold: 1 << 63, heap: 1 << 63, want: true, heapReads: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalHeap, originalGoroutines := readHeapAlloc, numGoroutine
			t.Cleanup(func() {
				readHeapAlloc, numGoroutine = originalHeap, originalGoroutines
			})

			var heapReads, goroutineReads int
			readHeapAlloc = func() uint64 {
				heapReads++
				return tt.heap
			}
			numGoroutine = func() int {
				goroutineReads++
				return tt.goroutines
			}

			c := Config{
				CpuThreshold:       tt.cpuThreshold,
				HeapThreshold:      tt.heapThreshold,
				GoroutineThreshold: tt.goroutineLimit,
			}
			assert.Equal(t, tt.want, isOverloaded(c))
			assert.Equal(t, tt.heapReads, heapReads)
			assert.Equal(t, tt.goroutineReads, goroutineReads)
		})
	}
}

func TestReadHeapAlloc(t *testing.T) {
	// Exercise the real runtime metric, so a misspelled or unsupported metric
	// cannot silently leave the heap trigger disabled.
	assert.Positive(t, readHeapAlloc())
}

func TestStartPyroscopeResourceTriggers(t *testing.T) {
	tests := []struct {
		name       string
		heap       uint64
		goroutines int64
	}{
		{name: "heap while CPU is below threshold", heap: 100},
		{name: "goroutines while CPU is below threshold", goroutines: 10},
		{name: "simultaneous resource thresholds", heap: 100, goroutines: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalFactory := newProfiler
			originalHeap, originalGoroutines := readHeapAlloc, numGoroutine
			t.Cleanup(func() {
				newProfiler = originalFactory
				readHeapAlloc, numGoroutine = originalHeap, originalGoroutines
			})

			var overloaded atomic.Bool
			var checks, creations atomic.Int32
			readHeapAlloc = func() uint64 {
				checks.Add(1)
				if overloaded.Load() {
					return tt.heap
				}
				return 0
			}
			numGoroutine = func() int {
				if overloaded.Load() {
					return int(tt.goroutines)
				}
				return 0
			}
			mp := &mockProfiler{}
			newProfiler = func(Config) profiler {
				creations.Add(1)
				// Recover in the profiling loop itself, so a delayed test
				// goroutine cannot accidentally permit a second session.
				overloaded.Store(false)
				return mp
			}

			// Use an unreachable CPU threshold to isolate the resource triggers
			// from the host load. ProfileType flags deliberately remain false:
			// choosing the collected profile types must not disable a trigger.
			c := Config{
				CpuThreshold:       math.MaxInt64,
				HeapThreshold:      100,
				GoroutineThreshold: 10,
				CheckInterval:      time.Millisecond,
				ProfilingDuration:  30 * time.Millisecond,
			}
			stop := startPyroscopeForTest(t, c)

			if !assert.Eventually(t, func() bool { return checks.Load() >= 3 },
				time.Second, time.Millisecond) {
				return
			}
			assert.Zero(t, creations.Load(), "below-threshold samples must not start a profile")

			overloaded.Store(true)
			if !assert.Eventually(t, mp.started.True, time.Second, time.Millisecond) {
				return
			}

			assert.Eventually(t, mp.stopped.True, time.Second, time.Millisecond)
			checksAfterStop := checks.Load()
			assert.Eventually(t, func() bool { return checks.Load() >= checksAfterStop+3 },
				time.Second, time.Millisecond)
			stop()
			assert.EqualValues(t, 1, creations.Load(), "only one profile should be started")
		})
	}
}
