package devserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_fillDefault(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected Config
	}{
		{
			"empty config should filled",
			Config{},
			Config{Port: 6470, EnableMetrics: true, MetricsPath: "/metrics", HealthPath: "/healthz"},
		},
		{
			"non empty config should not filled",
			Config{EnablePprof: true},
			Config{EnablePprof: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.config)
		})
	}
}
