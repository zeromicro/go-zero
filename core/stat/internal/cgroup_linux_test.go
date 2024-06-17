package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunningInUserNS(t *testing.T) {
	// should be false in docker
	assert.False(t, runningInUserNS())
}

func TestCgroups(t *testing.T) {
	// test cgroup legacy(v1) & hybrid
	if !isCgroup2UnifiedMode() {
		cg, err := currentCgroupV1()
		assert.NoError(t, err)
		_, err = cg.effectiveCpus()
		assert.NoError(t, err)
		_, err = cg.cpuQuota()
		assert.NoError(t, err)
		_, err = cg.cpuUsage()
		assert.NoError(t, err)
	}

	// test cgroup v2
	if isCgroup2UnifiedMode() {
		cg, err := currentCgroupV2()
		assert.NoError(t, err)
		_, err = cg.effectiveCpus()
		assert.NoError(t, err)
		_, err = cg.cpuQuota()
		assert.Error(t, err)
		_, err = cg.cpuUsage()
		assert.NoError(t, err)
	}
}

func TestParseUint(t *testing.T) {
	tests := []struct {
		input string
		want  uint64
		err   error
	}{
		{"0", 0, nil},
		{"123", 123, nil},
		{"-1", 0, nil},
		{"-18446744073709551616", 0, nil},
		{"foo", 0, fmt.Errorf("cgroup: bad int format: foo")},
	}

	for _, tt := range tests {
		got, err := parseUint(tt.input)
		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.want, got)
	}
}

func TestParseUints(t *testing.T) {
	tests := []struct {
		input string
		want  []uint64
		err   error
	}{
		{"", nil, nil},
		{"1,2,3", []uint64{1, 2, 3}, nil},
		{"1-3", []uint64{1, 2, 3}, nil},
		{"1-3,5,7-9", []uint64{1, 2, 3, 5, 7, 8, 9}, nil},
		{"foo", nil, fmt.Errorf("cgroup: bad int format: foo")},
		{"1-bar", nil, fmt.Errorf("cgroup: bad int list format: 1-bar")},
		{"bar-3", nil, fmt.Errorf("cgroup: bad int list format: bar-3")},
		{"3-1", nil, fmt.Errorf("cgroup: bad int list format: 3-1")},
	}

	for _, tt := range tests {
		got, err := parseUints(tt.input)
		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.want, got)
	}
}
