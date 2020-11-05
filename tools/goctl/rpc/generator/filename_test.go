package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatFilename(t *testing.T) {
	assert.Equal(t, "abc", formatFilename("a_b_c", namingLower))
	assert.Equal(t, "ABC", formatFilename("a_b_c", namingCamel))
	assert.Equal(t, "a_b_c", formatFilename("a_b_c", namingSnake))
	assert.Equal(t, "a", formatFilename("a", namingSnake))
	assert.Equal(t, "A", formatFilename("a", namingCamel))
	// no flag to convert to snake
	assert.Equal(t, "abc", formatFilename("abc", namingSnake))
}
