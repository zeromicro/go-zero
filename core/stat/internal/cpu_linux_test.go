package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefreshCpu(t *testing.T) {
	assert.True(t, RefreshCpu() >= 0)
}

func BenchmarkRefreshCpu(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RefreshCpu()
	}
}
