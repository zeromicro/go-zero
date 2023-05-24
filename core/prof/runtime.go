package prof

import (
	"fmt"
	"io"
	"os"
	"runtime"
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
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Fprintf(writer, "Goroutines: %d, Alloc: %vm, TotalAlloc: %vm, Sys: %vm, NumGC: %v\n",
				runtime.NumGoroutine(), m.Alloc/mega, m.TotalAlloc/mega, m.Sys/mega, m.NumGC)
		}
	}()
}
