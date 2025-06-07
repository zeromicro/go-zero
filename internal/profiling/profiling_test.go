package profiling

import (
	"sync"
	"testing"
	"time"

	"github.com/grafana/pyroscope-go"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestStart(t *testing.T) {
	t.Run("profiling", func(t *testing.T) {
		var c Config
		assert.NoError(t, conf.FillDefault(&c))
		c.Name = "test"
		p := newProfiler(c)
		assert.NotNil(t, p)
		assert.NoError(t, p.Start())
		assert.NoError(t, p.Stop())
	})

	t.Run("invalid config", func(t *testing.T) {
		mp := &mockProfiler{}
		newProfiler = func(c Config) profiler {
			return mp
		}

		Start(Config{})

		Start(Config{
			ServerAddr: "localhost:4040",
		})
	})

	t.Run("test start profiler", func(t *testing.T) {
		mp := &mockProfiler{}
		newProfiler = func(c Config) profiler {
			return mp
		}

		c := Config{
			Name:              "test",
			ServerAddr:        "localhost:4040",
			CheckInterval:     time.Millisecond,
			ProfilingDuration: time.Millisecond * 10,
			CpuThreshold:      0,
		}
		var done = make(chan struct{})
		go startPyroscope(c, done)

		time.Sleep(time.Millisecond * 50)
		close(done)

		assert.True(t, mp.started)
		assert.True(t, mp.stopped)
	})

	t.Run("test start profiler with cpu overloaded", func(t *testing.T) {
		mp := &mockProfiler{}
		newProfiler = func(c Config) profiler {
			return mp
		}

		c := Config{
			Name:              "test",
			ServerAddr:        "localhost:4040",
			CheckInterval:     time.Millisecond,
			ProfilingDuration: time.Millisecond * 10,
			CpuThreshold:      900,
		}
		var done = make(chan struct{})
		go startPyroscope(c, done)

		time.Sleep(time.Millisecond * 50)
		close(done)

		assert.False(t, mp.started)
	})

	t.Run("start/stop err", func(t *testing.T) {
		mp := &mockProfiler{
			err: assert.AnError,
		}
		newProfiler = func(c Config) profiler {
			return mp
		}

		c := Config{
			Name:              "test",
			ServerAddr:        "localhost:4040",
			CheckInterval:     time.Millisecond,
			ProfilingDuration: time.Millisecond * 10,
			CpuThreshold:      0,
		}
		var done = make(chan struct{})
		go startPyroscope(c, done)

		time.Sleep(time.Millisecond * 50)
		close(done)

		assert.False(t, mp.started)
		assert.False(t, mp.stopped)
	})
}

func TestGenPyroscopeConf(t *testing.T) {
	c := Config{
		Name:         "",
		ServerAddr:   "localhost:4040",
		AuthUser:     "user",
		AuthPassword: "password",
		ProfileType: ProfileType{
			Logger:     true,
			CPU:        true,
			Goroutines: true,
			Memory:     true,
			Mutex:      true,
			Block:      true,
		},
	}

	conf := genPyroscopeConf(c)
	assert.Equal(t, c.ServerAddr, conf.ServerAddress)
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

	newPyroscopeProfiler(c)
}

func TestNewPyroscopeProfiler(t *testing.T) {
	p := newPyroscopeProfiler(Config{})

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
