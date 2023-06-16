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

func TestCgroupV1(t *testing.T) {
	if isCgroup2UnifiedMode() {
		cg, err := currentCgroupV1()
		assert.NoError(t, err)
		_, err = cg.cpus()
		assert.Error(t, err)
		_, err = cg.cpuPeriodUs()
		assert.Error(t, err)
		_, err = cg.cpuQuotaUs()
		assert.Error(t, err)
		_, err = cg.usageAllCpus()
		assert.Error(t, err)
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
		{"foo", nil, fmt.Errorf("strconv.ParseUint: parsing \"foo\": invalid syntax")},
		{"1-bar", nil, fmt.Errorf("strconv.ParseUint: parsing \"1-bar\": invalid syntax")},
		{"3-1", nil, fmt.Errorf("cgroup: bad int list format: 3-1")},
	}

	for _, tt := range tests {
		got, err := parseUints(tt.input)
		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.want, got)
	}
}
