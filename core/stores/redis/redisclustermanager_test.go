package redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	red "github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
	"github.com/stretchr/testify/assert"
)

func TestSplitClusterAddrs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty input",
			input:    "",
			expected: []string{""},
		},
		{
			name:     "single address",
			input:    "127.0.0.1:8000",
			expected: []string{"127.0.0.1:8000"},
		},
		{
			name:     "multiple addresses with duplicates",
			input:    "127.0.0.1:8000,127.0.0.1:8001, 127.0.0.1:8000",
			expected: []string{"127.0.0.1:8000", "127.0.0.1:8001"},
		},
		{
			name:     "multiple addresses without duplicates",
			input:    "127.0.0.1:8000, 127.0.0.1:8001, 127.0.0.1:8002",
			expected: []string{"127.0.0.1:8000", "127.0.0.1:8001", "127.0.0.1:8002"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert.ElementsMatch(t, tc.expected, splitClusterAddrs(tc.input))
		})
	}
}

func TestGetCluster(t *testing.T) {
	r := miniredis.RunT(t)
	defer r.Close()
	c, err := getCluster(&Redis{
		Addr:  r.Addr(),
		Type:  ClusterType,
		tls:   true,
		hooks: []red.Hook{defaultDurationHook},
	})
	if assert.NoError(t, err) {
		assert.NotNil(t, c)
	}
}

func TestGetClusterWithProtocolAndIdentity(t *testing.T) {
	r := miniredis.RunT(t)
	defer r.Close()
	c, err := getCluster(&Redis{
		Addr:     r.Addr(),
		Type:     ClusterType,
		protocol: 2,
		identity: true,
		hooks:    []red.Hook{defaultDurationHook},
	})
	if assert.NoError(t, err) {
		assert.NotNil(t, c)
		assert.Equal(t, 2, c.Options().Protocol)
		assert.True(t, c.Options().DisableIdentity)
	}
}

func TestGetClusterWithMaintNotifications(t *testing.T) {
	tests := []struct {
		name string
		mode maintnotifications.Mode
		want maintnotifications.Mode
	}{
		{name: "unset falls back to disabled", mode: "", want: maintnotifications.ModeDisabled},
		{name: "disabled", mode: maintnotifications.ModeDisabled, want: maintnotifications.ModeDisabled},
		{name: "auto", mode: maintnotifications.ModeAuto, want: maintnotifications.ModeAuto},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := miniredis.RunT(t)
			defer r.Close()
			c, err := getCluster(&Redis{
				Addr:               r.Addr(),
				Type:               ClusterType,
				maintNotifications: test.mode,
				hooks:              []red.Hook{defaultDurationHook},
			})
			if assert.NoError(t, err) {
				assert.NotNil(t, c)
				assert.NotNil(t, c.Options().MaintNotificationsConfig)
				assert.Equal(t, test.want, c.Options().MaintNotificationsConfig.Mode)
			}
		})
	}
}
