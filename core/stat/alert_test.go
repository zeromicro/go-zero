//go:build linux

package stat

import (
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReport(t *testing.T) {
	t.Setenv(clusterNameKey, "test-cluster")

	var count int32
	SetReporter(func(s string) {
		atomic.AddInt32(&count, 1)
	})
	for i := 0; i < 10; i++ {
		Report(strconv.Itoa(i))
	}
	assert.Equal(t, int32(1), count)
}
