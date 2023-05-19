package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunningInUserNS(t *testing.T) {
	// should be false in docker
	assert.False(t, runningInUserNS())
}

func TestCgroupV2(t *testing.T) {
	cg, err := currentCgroup()
	assert.NoError(t, err)
	val, err := cg.cpuPeriodUs()
	assert.NoError(t, err)
	assert.NotEqual(t, 0, val)
}

func TestCgroupV1(t *testing.T) {
	if isCgroup2UnifiedMode() {
		_, err := currentCgroupV1()
		assert.Error(t, err)
	}
}
