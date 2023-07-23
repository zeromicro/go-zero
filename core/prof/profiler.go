package prof

import "github.com/zeromicro/go-zero/core/utils"

type (
	// A ProfilePoint is a profile time point.
	ProfilePoint struct {
		*utils.ElapsedTimer
	}

	// A Profiler interface represents a profiler that used to report profile points.
	Profiler interface {
		Start() ProfilePoint
		Report(name string, point ProfilePoint)
	}

	realProfiler struct{}

	nullProfiler struct{}
)

var profiler = newNullProfiler()

// EnableProfiling enables profiling.
func EnableProfiling() {
	profiler = newRealProfiler()
}

// Start starts a Profiler, and returns a start profiling point.
func Start() ProfilePoint {
	return profiler.Start()
}

// Report reports a ProfilePoint with given name.
func Report(name string, point ProfilePoint) {
	profiler.Report(name, point)
}

func newRealProfiler() Profiler {
	return &realProfiler{}
}

func (rp *realProfiler) Start() ProfilePoint {
	return ProfilePoint{
		ElapsedTimer: utils.NewElapsedTimer(),
	}
}

func (rp *realProfiler) Report(name string, point ProfilePoint) {
	duration := point.Duration()
	report(name, duration)
}

func newNullProfiler() Profiler {
	return &nullProfiler{}
}

func (np *nullProfiler) Start() ProfilePoint {
	return ProfilePoint{}
}

func (np *nullProfiler) Report(string, ProfilePoint) {
}
