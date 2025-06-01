package profiling

import (
	"sync"
	"testing"
	"time"

	"github.com/grafana/pyroscope-go"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	t.Run("invalid config", func(t *testing.T) {
		var mockProfiler = &mockProfiler{}
		newProfiler = func(c Config) profiler {
			return mockProfiler
		}

		Start(Config{})

		Start(Config{
			ServerAddress: "localhost:4040",
		})
	})

	t.Run("test start profiler", func(t *testing.T) {
		var mockProfiler = &mockProfiler{}
		newProfiler = func(c Config) profiler {
			return mockProfiler
		}

		c := Config{
			Name:              "test",
			ServerAddress:     "localhost:4040",
			IntervalDuration:  time.Millisecond,
			ProfilingDuration: time.Millisecond * 10,
			CpuThreshold:      0,
		}
		var done = make(chan struct{})
		go startPyroScope(c, done)

		time.Sleep(time.Millisecond * 50)
		done <- struct{}{}

		assert.True(t, mockProfiler.started)
		assert.True(t, mockProfiler.stopped)
	})

	t.Run("start/stop err", func(t *testing.T) {
		var mockProfiler = &mockProfiler{
			err: assert.AnError,
		}
		newProfiler = func(c Config) profiler {
			return mockProfiler
		}

		c := Config{
			Name:              "test",
			ServerAddress:     "localhost:4040",
			IntervalDuration:  time.Millisecond,
			ProfilingDuration: time.Millisecond * 10,
			CpuThreshold:      0,
		}
		var done = make(chan struct{})
		go startPyroScope(c, done)

		time.Sleep(time.Millisecond * 50)
		done <- struct{}{}

		assert.False(t, mockProfiler.started)
		assert.False(t, mockProfiler.stopped)
	})
}

func TestGenPyroScopeConf(t *testing.T) {
	c := Config{
		Name:          "",
		ServerAddress: "localhost:4040",
		AuthUser:      "user",
		AuthPassword:  "password",
		ProfileType: ProfileType{
			Logger:     true,
			CPU:        true,
			Goroutines: true,
			Memory:     true,
			Mutex:      true,
			Block:      true,
		},
	}

	conf := genPyroScopeConf(c)
	assert.Equal(t, c.ServerAddress, conf.ServerAddress)
	assert.Equal(t, c.AuthUser, conf.BasicAuthUser)
	assert.Equal(t, c.AuthPassword, conf.BasicAuthPassword)
	assert.Equal(t, c.Name, conf.ApplicationName)
	assert.Contains(t, conf.ProfileTypes, pyroscope.ProfileCPU)
	assert.Contains(t, conf.ProfileTypes, pyroscope.ProfileGoroutines)
	assert.Contains(t, conf.ProfileTypes, pyroscope.ProfileAllocObjects)
	assert.Contains(t, conf.ProfileTypes, pyroscope.ProfileAllocSpace)
	assert.Contains(t, conf.ProfileTypes, pyroscope.ProfileInuseObjects)
	assert.Contains(t, conf.ProfileTypes, pyroscope.ProfileInuseSpace)
	assert.Contains(t, conf.ProfileTypes, pyroscope.ProfileMutexCount)
	assert.Contains(t, conf.ProfileTypes, pyroscope.ProfileMutexDuration)
	assert.Contains(t, conf.ProfileTypes, pyroscope.ProfileBlockCount)
	assert.Contains(t, conf.ProfileTypes, pyroscope.ProfileBlockDuration)

	setFraction(c)
	resetFraction(c)

	newPyProfiler(c)
}

func TestNewPyProfiler(t *testing.T) {
	p := newPyProfiler(Config{})

	assert.Error(t, p.Start())
	assert.NoError(t, p.Stop())
}

type mockProfiler struct {
	mutex   sync.Mutex
	started bool
	stopped bool
	err     error
}

func (m *mockProfiler) Start() error {
	m.mutex.Lock()
	if m.err == nil {
		m.started = true
	}
	m.mutex.Unlock()
	return m.err
}

func (m *mockProfiler) Stop() error {
	m.mutex.Lock()
	if m.err == nil {
		m.stopped = true
	}
	m.mutex.Unlock()
	return m.err
}
