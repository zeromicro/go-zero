package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/internal/devserver"
)

func TestServiceConfProfilingThresholds(t *testing.T) {
	tests := []struct {
		name       string
		settings   string
		cpu        int64
		heap       uint64
		goroutines int
		wantErr    bool
	}{
		{name: "existing configuration", cpu: 700},
		{name: "independent resource thresholds", cpu: 700, heap: 536870912, goroutines: 10000,
			settings: "  HeapThreshold: 536870912\n  GoroutineThreshold: 10000\n"},
		{name: "explicit zero thresholds", settings: "  CpuThreshold: 0\n  HeapThreshold: 0\n  GoroutineThreshold: 0\n"},
		{name: "negative heap", settings: "  HeapThreshold: -1\n", wantErr: true},
		{name: "negative goroutines", settings: "  GoroutineThreshold: -1\n", wantErr: true},
		{name: "fractional heap", settings: "  HeapThreshold: 1.5\n", wantErr: true},
		{name: "fractional goroutines", settings: "  GoroutineThreshold: 1.5\n", wantErr: true},
		{name: "CPU upper bound unchanged", settings: "  CpuThreshold: 1000\n", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c ServiceConf
			yaml := fmt.Sprintf("Name: worker\nProfiling:\n  ServerAddr: http://localhost:4040\n%s", tt.settings)
			err := conf.LoadFromYamlBytes([]byte(yaml), &c)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tt.cpu, c.Profiling.CpuThreshold)
			assert.Equal(t, tt.heap, c.Profiling.HeapThreshold)
			assert.Equal(t, tt.goroutines, c.Profiling.GoroutineThreshold)
			assert.Equal(t, "worker", c.Profiling.Name)
			assert.True(t, c.Profiling.ProfileType.Memory)
			assert.True(t, c.Profiling.ProfileType.Goroutines)
		})
	}
}

func TestServiceConf(t *testing.T) {
	c := ServiceConf{
		Name: "foo",
		Log: logx.LogConf{
			Mode: "console",
		},
		Mode: "dev",
		DevServer: devserver.Config{
			Port:       6470,
			HealthPath: "/healthz",
		},
	}
	c.MustSetUp()
}

func TestServiceConfWithMetricsUrl(t *testing.T) {
	c := ServiceConf{
		Name: "foo",
		Log: logx.LogConf{
			Mode: "volume",
		},
		Mode:       "dev",
		MetricsUrl: "http://localhost:8080",
	}
	assert.NoError(t, c.SetUp())
}
