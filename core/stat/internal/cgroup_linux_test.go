package internal

import (
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
