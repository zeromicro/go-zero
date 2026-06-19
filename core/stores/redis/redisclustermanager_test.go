package redis

import (
	"crypto/tls"
	"testing"

	"github.com/alicebob/miniredis/v2"
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
		Addr:      r.Addr(),
		Type:      ClusterType,
		tlsConfig: &tls.Config{},
	})
	if assert.NoError(t, err) {
		assert.NotNil(t, c)
		assert.False(t, c.Options().TLSConfig.InsecureSkipVerify)
	}
}
