package stat

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx/logtest"
)

func TestBToMb(t *testing.T) {
	tests := []struct {
		name     string
		bytes    uint64
		expected float32
	}{
		{
			name:     "Test 1: Convert 0 bytes to MB",
			bytes:    0,
			expected: 0,
		},
		{
			name:     "Test 2: Convert 1048576 bytes to MB",
			bytes:    1048576,
			expected: 1,
		},
		{
			name:     "Test 3: Convert 2097152 bytes to MB",
			bytes:    2097152,
			expected: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := bToMb(test.bytes)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestPrintUsage(t *testing.T) {
	c := logtest.NewCollector(t)

	printUsage()

	output := c.String()
	assert.Contains(t, output, "CPU:")
	assert.Contains(t, output, "MEMORY:")
	assert.Contains(t, output, "Alloc=")
	assert.Contains(t, output, "TotalAlloc=")
	assert.Contains(t, output, "Sys=")
	assert.Contains(t, output, "NumGC=")

	lines := strings.Split(output, "\n")
	assert.Len(t, lines, 2)
	fields := strings.Split(lines[0], ", ")
	assert.Len(t, fields, 5)
}
