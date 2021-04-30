package name

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNamingValid(t *testing.T) {
	style, valid := IsNamingValid("")
	assert.True(t, valid)
	assert.Equal(t, NamingLower, style)

	_, valid = IsNamingValid("lower1")
	assert.False(t, valid)

	_, valid = IsNamingValid("lower")
	assert.True(t, valid)

	_, valid = IsNamingValid("snake")
	assert.True(t, valid)

	_, valid = IsNamingValid("camel")
	assert.True(t, valid)
}

func TestFormatFilename(t *testing.T) {
	assert.Equal(t, "abc", FormatFilename("a_b_c", NamingLower))
	assert.Equal(t, "ABC", FormatFilename("a_b_c", NamingCamel))
	assert.Equal(t, "a_b_c", FormatFilename("a_b_c", NamingSnake))
	assert.Equal(t, "a", FormatFilename("a", NamingSnake))
	assert.Equal(t, "A", FormatFilename("a", NamingCamel))
	// no flag to convert to snake
	assert.Equal(t, "abc", FormatFilename("abc", NamingSnake))
}
