package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatKeys(t *testing.T) {
	assert.Equal(t, "a,b", formatKeys([]string{"a", "b"}))
}

func TestTotalWeights(t *testing.T) {
	val := TotalWeights([]NodeConf{
		{
			Weight: -1,
		},
		{
			Weight: 0,
		},
		{
			Weight: 1,
		},
	})
	assert.Equal(t, 1, val)
}
