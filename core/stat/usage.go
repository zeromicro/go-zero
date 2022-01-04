package stat

import (
	"runtime"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stat/internal"
	"github.com/zeromicro/go-zero/core/threading"
)

const (
	// 250ms and 0.95 as beta will count the average cpu load for past 5 seconds
	cpuRefreshInterval = time.Millisecond * 250
	allRefreshInterval = time.Minute
	// moving average beta hyperparameter
	beta = 0.95
)

var cpuUsage int64

func init() {
	go func() {
		cpuTicker := time.NewTicker(cpuRefreshInterval)
		defer cpuTicker.Stop()
		allTicker := time.NewTicker(allRefreshInterval)
		defer allTicker.Stop()

		for {
			select {
			case <-cpuTicker.C:
				threading.RunSafe(func() {
					curUsage := internal.RefreshCpu()
					prevUsage := atomic.LoadInt64(&cpuUsage)
					// cpu = cpuᵗ⁻¹ * beta + cpuᵗ * (1 - beta)
					usage := int64(float64(prevUsage)*beta + float64(curUsage)*(1-beta))
					atomic.StoreInt64(&cpuUsage, usage)
				})
			case <-allTicker.C:
				if logEnabled.True() {
					printUsage()
				}
			}
		}
	}()
}

// CpuUsage returns current cpu usage.
func CpuUsage() int64 {
	return atomic.LoadInt64(&cpuUsage)
}

func bToMb(b uint64) float32 {
	return float32(b) / 1024 / 1024
}

func printUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	logx.Statf("CPU: %dm, MEMORY: Alloc=%.1fMi, TotalAlloc=%.1fMi, Sys=%.1fMi, NumGC=%d",
		CpuUsage(), bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
}
