package continuousprofiling

import (
	"bytes"
	"testing"
	"time"

	"github.com/grafana/pyroscope-go"
	"github.com/stretchr/testify/assert"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestStart(t *testing.T) {
	var buf bytes.Buffer
	logx.Reset()
	logx.SetLevel(logx.InfoLevel)
	logx.SetWriter(logx.NewWriter(&buf))
	defer logx.Reset()

	newProfiler = func(c Config) profiler {
		return &mockProfiler{}
	}

	c := Config{}
	conf.FillDefault(&c)

	c.ServerAddress = "localhost:4040"
	c.IntervalDuration = time.Millisecond
	c.ProfilingDuration = time.Millisecond * 10
	c.CpuThreshold = 0

	Start(c)

	time.Sleep(time.Second * 1)

	assert.Contains(t, buf.String(), "pyroScope profiler started")
	assert.Contains(t, buf.String(), "pyroScope profiler stopped")
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
}

type mockProfiler struct {
	err error
}

func (m *mockProfiler) Start() error {

	return m.err
}

func (m *mockProfiler) Stop() error {
	return m.err
}
