package continueprofiling

import (
	"testing"

	"github.com/grafana/pyroscope-go"
	"github.com/stretchr/testify/assert"
)

func TestGenPyroScopeConf(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		contains    []pyroscope.ProfileType
		notContains []pyroscope.ProfileType
	}{
		{
			name: "default profile types",
			config: Config{
				ProfileType: ProfileType{
					CPUOff:        false,
					GoroutinesOff: false,
					MemoryOff:     false,
					MutexOff:      true,
					BlockOff:      true,
				},
			},
			contains: []pyroscope.ProfileType{
				pyroscope.ProfileCPU,
				pyroscope.ProfileGoroutines,
				pyroscope.ProfileAllocObjects,
				pyroscope.ProfileAllocSpace,
				pyroscope.ProfileInuseObjects,
				pyroscope.ProfileInuseSpace,
			},
			notContains: []pyroscope.ProfileType{
				pyroscope.ProfileMutexCount,
				pyroscope.ProfileMutexDuration,
				pyroscope.ProfileBlockCount,
				pyroscope.ProfileBlockDuration,
			},
		},
		{
			name: "cpu profiling off",
			config: Config{
				ProfileType: ProfileType{
					CPUOff: true,
				},
			},
			notContains: []pyroscope.ProfileType{pyroscope.ProfileCPU},
		},
		{
			name: "auth credentials",
			config: Config{
				AuthUser:     "testuser",
				AuthPassword: "testpass",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := genPyroScopeConf(tt.config)

			// Verify contained profile types
			for _, pt := range tt.contains {
				assert.Contains(t, conf.ProfileTypes, pt)
			}

			// Verify excluded profile types
			for _, pt := range tt.notContains {
				assert.NotContains(t, conf.ProfileTypes, pt)
			}

			// Verify auth credentials
			if tt.config.AuthUser != "" {
				assert.Equal(t, tt.config.AuthUser, conf.BasicAuthUser)
				assert.Equal(t, tt.config.AuthPassword, conf.BasicAuthPassword)
			}
		})
	}
}
