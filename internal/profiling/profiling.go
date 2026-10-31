package profiling

import (
	"runtime"
	"runtime/metrics"
	"sync"
	"time"

	"github.com/grafana/pyroscope-go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/threading"
)

const (
	defaultCheckInterval     = time.Second * 10
	defaultProfilingDuration = time.Minute * 2
	defaultUploadRate        = time.Second * 15
)

type (
	Config struct {
		// Name is the name of the application.
		Name string `json:",optional,inherit"`
		// ServerAddr is the address of the profiling server.
		ServerAddr string
		// AuthUser is the username for basic authentication.
		AuthUser string `json:",optional"`
		// AuthPassword is the password for basic authentication.
		AuthPassword string `json:",optional"`
		// UploadRate is the duration for which profiling data is uploaded.
		UploadRate time.Duration `json:",default=15s"`
		// CheckInterval is the interval to check if profiling should start.
		CheckInterval time.Duration `json:",default=10s"`
		// ProfilingDuration is the duration for which profiling data is collected.
		ProfilingDuration time.Duration `json:",default=2m"`
		// CpuThreshold is the CPU usage threshold on a scale of 0 to 1000 (1000 = 100%).
		// A value of 0 triggers profiling unconditionally.
		CpuThreshold int64 `json:",default=700,range=[0:1000)"`
		// HeapThreshold is the allocated heap memory threshold in bytes.
		// A value of 0 disables the heap trigger.
		HeapThreshold uint64 `json:",default=0"`
		// GoroutineThreshold is the threshold for the number of live goroutines.
		// A value of 0 disables the goroutine trigger.
		GoroutineThreshold int `json:",default=0,range=[0:)"`

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

	pyroscopeProfiler struct {
		c        Config
		profiler *pyroscope.Profiler
	}
)

var (
	once sync.Once

	newProfiler = func(c Config) profiler {
		return newPyroscopeProfiler(c)
	}
	numGoroutine  = runtime.NumGoroutine
	readHeapAlloc = func() uint64 {
		samples := []metrics.Sample{
			{Name: "/memory/classes/heap/objects:bytes"},
		}
		metrics.Read(samples)
		if samples[0].Value.Kind() == metrics.KindUint64 {
			return samples[0].Value.Uint64()
		}

		return 0
	}
)

// Start initializes the pyroscope profiler with the given configuration.
func Start(c Config) {
	// check if the profiling is enabled
	if len(c.ServerAddr) == 0 {
		return
	}

	// set default values for the configuration
	if c.ProfilingDuration <= 0 {
		c.ProfilingDuration = defaultProfilingDuration
	}

	// set default values for the configuration
	if c.CheckInterval <= 0 {
		c.CheckInterval = defaultCheckInterval
	}

	if c.UploadRate <= 0 {
		c.UploadRate = defaultUploadRate
	}

	once.Do(func() {
		logx.Info("continuous profiling started")

		threading.GoSafe(func() {
			startPyroscope(c, proc.Done())
		})
	})
}

// startPyroscope starts the pyroscope profiler with the given configuration.
func startPyroscope(c Config, done <-chan struct{}) {
	var (
		pr                  profiler
		err                 error
		latestProfilingTime time.Time
		intervalTicker      = time.NewTicker(c.CheckInterval)
		profilingTicker     = time.NewTicker(c.ProfilingDuration)
	)

	defer profilingTicker.Stop()
	defer intervalTicker.Stop()

	for {
		select {
		case <-intervalTicker.C:
			// Check if the machine is overloaded and if the profiler is not running
			if pr == nil && isOverloaded(c) {
				pr = newProfiler(c)
				if err := pr.Start(); err != nil {
					logx.Errorf("failed to start profiler: %v", err)
					continue
				}

				// record the latest profiling time
				latestProfilingTime = time.Now()
				logx.Infof("pyroscope profiler started.")
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
				logx.Infof("pyroscope profiler stopped.")
				pr = nil
			}
		case <-done:
			logx.Infof("continuous profiling stopped.")
			return
		}
	}
}

// genPyroscopeConf generates the pyroscope configuration based on the given config.
func genPyroscopeConf(c Config) pyroscope.Config {
	pConf := pyroscope.Config{
		UploadRate:        c.UploadRate,
		ApplicationName:   c.Name,
		BasicAuthUser:     c.AuthUser,     // http basic auth user
		BasicAuthPassword: c.AuthPassword, // http basic auth password
		ServerAddress:     c.ServerAddr,
		Logger:            nil,
		HTTPHeaders:       map[string]string{},
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

// isCpuOverloaded checks the machine performance based on the given configuration.
func isCpuOverloaded(c Config) bool {
	currentValue := stat.CpuUsage()
	if currentValue >= c.CpuThreshold {
		logx.Infof("continuous profiling cpu overload, cpu: %d", currentValue)
		return true
	}

	return false
}

// isOverloaded checks whether any enabled profiling threshold has been reached.
func isOverloaded(c Config) bool {
	if isCpuOverloaded(c) {
		return true
	}

	if c.HeapThreshold > 0 {
		currentValue := readHeapAlloc()
		if currentValue >= c.HeapThreshold {
			logx.Infof("continuous profiling heap overload, heap: %d bytes", currentValue)
			return true
		}
	}

	if c.GoroutineThreshold > 0 {
		currentValue := numGoroutine()
		if currentValue >= c.GoroutineThreshold {
			logx.Infof("continuous profiling goroutine overload, goroutines: %d", currentValue)
			return true
		}
	}

	return false
}

func newPyroscopeProfiler(c Config) profiler {
	return &pyroscopeProfiler{
		c: c,
	}
}

func (p *pyroscopeProfiler) Start() error {
	pConf := genPyroscopeConf(p.c)
	// set mutex and block profile rate
	setFraction(p.c)
	prof, err := pyroscope.Start(pConf)
	if err != nil {
		resetFraction(p.c)
		return err
	}

	p.profiler = prof
	return nil
}

func (p *pyroscopeProfiler) Stop() error {
	if p.profiler == nil {
		return nil
	}

	if err := p.profiler.Stop(); err != nil {
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
