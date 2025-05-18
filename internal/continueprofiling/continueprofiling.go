package continueprofiling

import (
	"sync"
	"time"

	"github.com/grafana/pyroscope-go"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/threading"
)

type (
	Config struct {
		// Name is the name of the application.
		Name string `json:",optional,inherit"`
		// ServerAddress is the address of the profiling server.
		ServerAddress string
		// AuthUser is the username for basic authentication.
		AuthUser string `json:",optional"`
		// AuthPassword is the password for basic authentication.
		AuthPassword string `json:",optional"`
		// UploadDuration is the duration for which profiling data is uploaded.
		UploadDuration time.Duration `json:",default=15s"`
		// IntervalTicker is the interval for which profiling data is collected.
		IntervalTicker time.Duration `json:",default=10s"`
		// ProfilingDuration is the duration for which profiling data is collected.
		ProfilingDuration time.Duration `json:",default=2m"`
		// CpuThreshold the collection is allowed only when the current service cpu < CpuThreshold
		CpuThreshold int64 `json:",default=700,range=[0:1000)"`

		// ProfileType is the type of profiling to be performed.
		ProfileType ProfileType
	}

	ProfileType struct {
		// LoggerOpen is a flag to enable or disable logger profiling.
		LoggerOpen bool `json:",default=false"`
		// DisableGCRunsOff is a flag to disable garbage collection runs.
		DisableGCRunsOff bool `json:",default=false"`
		// CPUOff is a flag to disable CPU profiling.
		CPUOff bool `json:",default=false"`
		// GoroutinesOff is a flag to disable goroutine profiling.
		GoroutinesOff bool `json:",default=false"`
		// MemoryOff is a flag to disable memory profiling.
		MemoryOff bool `json:",default=false"`
		// MutexOff is a flag to disable mutex profiling.
		MutexOff bool `json:",default=true"`
		// BlockOff is a flag to disable block profiling.
		BlockOff bool `json:",default=true"`
	}
)

var once sync.Once

// Start initializes the pyroscope profiler with the given configuration.
func Start(c Config) {
	// check if the profiling is enabled
	if c.ServerAddress == "" {
		return
	}

	// set default values for the configuration
	if c.ProfilingDuration <= 0 {
		c.ProfilingDuration = time.Minute * 2
	}

	// set default values for the configuration
	if c.IntervalTicker <= 0 {
		c.IntervalTicker = time.Second * 10
	}

	if c.UploadDuration <= 0 {
		c.UploadDuration = time.Second * 15
	}

	once.Do(func() {
		logx.Info("continuous profiling started")
		threading.GoSafe(func() {
			startPyroScope(c)
		})
	})
}

// startPyroScope starts the pyroscope profiler with the given configuration.
func startPyroScope(c Config) {
	var (
		intervalTicker  = time.NewTicker(c.IntervalTicker)
		profilingTicker = time.NewTicker(c.ProfilingDuration)

		profiler *pyroscope.Profiler
		err      error

		latestProfilingTime time.Time
	)

	for {
		select {
		case <-intervalTicker.C:
			// Check if the machine is overloaded and if the profiler is not running
			if profiler == nil && checkMachinePerformance(c) {
				pConf := genPyroScopeConf(c)
				profiler, err = pyroscope.Start(pConf)
				if err != nil {
					logx.Errorf("failed to start profiler: %v", err)
					continue
				}

				// record the latest profiling time
				latestProfilingTime = time.Now()
				logx.Infof("pyroScope profiler started.")
			}
		case <-profilingTicker.C:
			// check if the profiling duration has passed
			if !time.Now().After(latestProfilingTime.Add(c.ProfilingDuration)) {
				continue
			}

			// check if the profiler is already running, if so, skip
			if profiler != nil {
				if err = profiler.Stop(); err != nil {
					logx.Errorf("failed to stop profiler: %v", err)
				}
				logx.Infof("pyroScope profiler stopped.")
				profiler = nil
			}
		}
	}
}

// genPyroScopeConf generates the pyroscope configuration based on the given config.
func genPyroScopeConf(c Config) pyroscope.Config {
	pConf := pyroscope.Config{
		UploadRate:        c.UploadDuration,
		DisableGCRuns:     !c.ProfileType.DisableGCRunsOff, // disable GC runs
		ApplicationName:   c.Name,
		BasicAuthUser:     c.AuthUser,     // http basic auth user
		BasicAuthPassword: c.AuthPassword, // http basic auth password
		// replace this with the address of pyroscope server
		ServerAddress: c.ServerAddress,

		// you can disable logging by setting this to nil
		Logger: logx.WithCallerSkip(0),

		HTTPHeaders: map[string]string{},

		// you can provide static tags via a map:
		Tags: map[string]string{
			"name": c.Name,
		},
	}

	if !c.ProfileType.CPUOff {
		pConf.ProfileTypes = append(pConf.ProfileTypes, pyroscope.ProfileCPU)
	}
	if !c.ProfileType.GoroutinesOff {
		pConf.ProfileTypes = append(pConf.ProfileTypes, pyroscope.ProfileGoroutines)
	}
	if !c.ProfileType.MemoryOff {
		pConf.ProfileTypes = append(pConf.ProfileTypes, pyroscope.ProfileAllocObjects, pyroscope.ProfileAllocSpace, pyroscope.ProfileInuseObjects, pyroscope.ProfileInuseSpace)
	}
	if !c.ProfileType.MutexOff {
		pConf.ProfileTypes = append(pConf.ProfileTypes, pyroscope.ProfileMutexCount, pyroscope.ProfileMutexDuration)
	}
	if !c.ProfileType.BlockOff {
		pConf.ProfileTypes = append(pConf.ProfileTypes, pyroscope.ProfileBlockCount, pyroscope.ProfileBlockDuration)
	}
	logx.Infof("applicationName: %s", pConf.ApplicationName)

	return pConf
}

// checkMachinePerformance checks the machine performance based on the given configuration.
func checkMachinePerformance(c Config) bool {
	currentValue := stat.CpuUsage()
	// overload >= 700, 70%
	if currentValue >= c.CpuThreshold {
		logx.Infof("continuous profiling cpu overload, cpu:%d", currentValue)
		return true
	}

	return false
}
