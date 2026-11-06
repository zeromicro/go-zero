package profiling

import (
	"sync"
	"testing"
	"time"

	"github.com/grafana/pyroscope-go"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/syncx"
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

		assert.True(t, mp.started.True())
		assert.True(t, mp.stopped.True())
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

		assert.False(t, mp.started.True())
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

		assert.False(t, mp.started.True())
		assert.False(t, mp.stopped.True())
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

	pyroscopeConf := genPyroscopeConf(c)
	assert.Equal(t, c.ServerAddr, pyroscopeConf.ServerAddress)
	assert.Equal(t, c.AuthUser, pyroscopeConf.BasicAuthUser)
	assert.Equal(t, c.AuthPassword, pyroscopeConf.BasicAuthPassword)
	assert.Equal(t, c.Name, pyroscopeConf.ApplicationName)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileCPU)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileGoroutines)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileAllocObjects)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileAllocSpace)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileInuseObjects)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileInuseSpace)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileMutexCount)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileMutexDuration)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileBlockCount)
	assert.Contains(t, pyroscopeConf.ProfileTypes, pyroscope.ProfileBlockDuration)

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
	started syncx.AtomicBool
	stopped syncx.AtomicBool
	err     error
}

func (m *mockProfiler) Start() error {
	m.mutex.Lock()
	if m.err == nil {
		m.started.Set(true)
	}
	m.mutex.Unlock()
	return m.err
}

func (m *mockProfiler) Stop() error {
	m.mutex.Lock()
	if m.err == nil {
		m.stopped.Set(true)
	}
	m.mutex.Unlock()
	return m.err
}
