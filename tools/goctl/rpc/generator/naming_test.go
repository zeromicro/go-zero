package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNamingValid(t *testing.T) {
	style, valid := IsNamingValid("")
	assert.True(t, valid)
	assert.Equal(t, namingLower, style)

	_, valid = IsNamingValid("lower1")
	assert.False(t, valid)

	_, valid = IsNamingValid("lower")
	assert.True(t, valid)

	_, valid = IsNamingValid("snake")
	assert.True(t, valid)

	_, valid = IsNamingValid("camel")
	assert.True(t, valid)
}
