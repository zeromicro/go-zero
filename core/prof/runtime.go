package prof

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/metrics"
	"time"
)

const (
	defaultInterval = time.Second * 5
	mega            = 1024 * 1024
)

// DisplayStats prints the goroutine, memory, GC stats with given interval, default to 5 seconds.
func DisplayStats(interval ...time.Duration) {
	displayStatsWithWriter(os.Stdout, interval...)
}

func displayStatsWithWriter(writer io.Writer, interval ...time.Duration) {
	duration := defaultInterval
	for _, val := range interval {
		duration = val
	}

	go func() {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		for range ticker.C {
			var (
				alloc, totalAlloc, sys uint64
				samples                = []metrics.Sample{
					{Name: "/memory/classes/heap/objects:bytes"},
					{Name: "/gc/heap/allocs:bytes"},
					{Name: "/memory/classes/total:bytes"},
				}
			)
			metrics.Read(samples)

			if samples[0].Value.Kind() == metrics.KindUint64 {
				alloc = samples[0].Value.Uint64()
			}
			if samples[1].Value.Kind() == metrics.KindUint64 {
				totalAlloc = samples[1].Value.Uint64()
			}
			if samples[2].Value.Kind() == metrics.KindUint64 {
				sys = samples[2].Value.Uint64()
			}
			var stats debug.GCStats
			debug.ReadGCStats(&stats)
			fmt.Fprintf(writer, "Goroutines: %d, Alloc: %vm, TotalAlloc: %vm, Sys: %vm, NumGC: %v\n",
				runtime.NumGoroutine(), alloc/mega, totalAlloc/mega, sys/mega, stats.NumGC)
		}
	}()
}
