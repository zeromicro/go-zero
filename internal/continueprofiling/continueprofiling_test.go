package continueprofiling

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestStartPyroScope(t *testing.T) {
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

	c.IntervalDuration = time.Millisecond
	c.ProfilingDuration = time.Millisecond * 10
	c.CpuThreshold = 0

	var done = make(chan struct{})

	go startPyroScope(c, done)

	time.Sleep(time.Second * 1)
	done <- struct{}{}

	assert.Contains(t, buf.String(), "pyroScope profiler started")
	assert.Contains(t, buf.String(), "pyroScope profiler stopped")
	assert.Contains(t, buf.String(), "continuous profiling stopped")
}

type mockProfiler struct {
}

func (m *mockProfiler) Start() error {

	return nil
}

func (m *mockProfiler) Stop() error {
	return nil
}
