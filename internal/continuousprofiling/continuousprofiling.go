package continuousprofiling

import (
	"runtime"
	"sync"
	"time"

	"github.com/grafana/pyroscope-go"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
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
		// IntervalDuration is the interval for which profiling data is collected.
		IntervalDuration time.Duration `json:",default=10s"`
		// ProfilingDuration is the duration for which profiling data is collected.
		ProfilingDuration time.Duration `json:",default=2m"`
		// CpuThreshold the collection is allowed only when the current service cpu < CpuThreshold
		CpuThreshold int64 `json:",default=700,range=[0:1000)"`

		// ProfileType is the type of profiling to be performed.
		ProfileType ProfileType
	}

	ProfileType struct {
		// Logger is a flag to enable or disable logging.
		Logger bool `json:",default=false"`
		// CPU is a flag to disable CPU profiling.
		CPU bool `json:",default=true"`
		// Goroutines is a flag to disable goroutine profiling.
		Goroutines bool `json:",default=true"`
		// Memory is a flag to disable memory profiling.
		Memory bool `json:",default=true"`
		// Mutex is a flag to disable mutex profiling.
		Mutex bool `json:",default=false"`
		// Block is a flag to disable block profiling.
		Block bool `json:",default=false"`
	}

	profiler interface {
		Start() error
		Stop() error
	}

	pyProfiler struct {
		c        Config
		profiler *pyroscope.Profiler
	}
)

var (
	once sync.Once

	newProfiler = func(c Config) profiler {
		return newPyProfiler(c)
	}
)

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
	if c.IntervalDuration <= 0 {
		c.IntervalDuration = time.Second * 10
	}

	if c.UploadDuration <= 0 {
		c.UploadDuration = time.Second * 15
	}

	once.Do(func() {
		logx.Info("continuous profiling started")
		var done = make(chan struct{})
		proc.AddShutdownListener(func() {
			done <- struct{}{}
		})

		threading.GoSafe(func() {
			startPyroScope(c, done)
		})
	})
}

// startPyroScope starts the pyroscope profiler with the given configuration.
func startPyroScope(c Config, done <-chan struct{}) {
	var (
		intervalTicker  = time.NewTicker(c.IntervalDuration)
		profilingTicker = time.NewTicker(c.ProfilingDuration)

		pr  profiler
		err error

		latestProfilingTime time.Time
	)

	for {
		select {
		case <-intervalTicker.C:
			// Check if the machine is overloaded and if the profiler is not running
			if pr == nil && checkMachinePerformance(c) {
				pr = newProfiler(c)
				if err := pr.Start(); err != nil {
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
			if pr != nil {
				if err = pr.Stop(); err != nil {
					logx.Errorf("failed to stop profiler: %v", err)
				}
				logx.Infof("pyroScope profiler stopped.")
				pr = nil
			}
		case <-done:
			logx.Infof("continuous profiling stopped.")
			return
		}
	}
}

// genPyroScopeConf generates the pyroscope configuration based on the given config.
func genPyroScopeConf(c Config) pyroscope.Config {
	pConf := pyroscope.Config{
		UploadRate:        c.UploadDuration,
		ApplicationName:   c.Name,
		BasicAuthUser:     c.AuthUser,     // http basic auth user
		BasicAuthPassword: c.AuthPassword, // http basic auth password
		ServerAddress:     c.ServerAddress,
		Logger:            nil,

		HTTPHeaders: map[string]string{},

		// you can provide static tags via a map:
		Tags: map[string]string{
			"name": c.Name,
		},
	}

	if c.ProfileType.Logger {
		pConf.Logger = logx.WithCallerSkip(0)
	}

	if c.ProfileType.CPU {
		pConf.ProfileTypes = append(pConf.ProfileTypes, pyroscope.ProfileCPU)
	}
	if c.ProfileType.Goroutines {
		pConf.ProfileTypes = append(pConf.ProfileTypes, pyroscope.ProfileGoroutines)
	}
	if c.ProfileType.Memory {
		pConf.ProfileTypes = append(pConf.ProfileTypes, pyroscope.ProfileAllocObjects, pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects, pyroscope.ProfileInuseSpace)
	}
	if c.ProfileType.Mutex {
		pConf.ProfileTypes = append(pConf.ProfileTypes, pyroscope.ProfileMutexCount, pyroscope.ProfileMutexDuration)
	}
	if c.ProfileType.Block {
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

func newPyProfiler(c Config) profiler {
	return &pyProfiler{
		c: c,
	}
}

func (p *pyProfiler) Start() error {
	pConf := genPyroScopeConf(p.c)
	// set mutex and block profile rate
	setFraction(p.c)
	profiler, err := pyroscope.Start(pConf)
	if err != nil {
		resetFraction(p.c)
		return err
	}

	p.profiler = profiler
	return nil
}

func (p *pyProfiler) Stop() error {
	if p.profiler == nil {
		return nil
	}

	err := p.profiler.Stop()
	if err != nil {
		return err
	}
	resetFraction(p.c)
	p.profiler = nil

	return nil
}

func setFraction(c Config) {
	// These 2 lines are only required if you're using mutex or block profiling
	if c.ProfileType.Mutex {
		runtime.SetMutexProfileFraction(10) // 10/seconds
	}
	if c.ProfileType.Block {
		runtime.SetBlockProfileRate(1000 * 1000) //  1/millisecond
	}
}

func resetFraction(c Config) {
	// These 2 lines are only required if you're using mutex or block profiling
	if c.ProfileType.Mutex {
		runtime.SetMutexProfileFraction(0)
	}
	if c.ProfileType.Block {
		runtime.SetBlockProfileRate(0)
	}
}
