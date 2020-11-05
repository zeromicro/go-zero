package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNamingValid(t *testing.T) {
	assert.True(t, func() bool {
		v, valid := IsNamingValid("lower")
		return v == namingLower && valid
	}())
	assert.True(t, func() bool {
		v, valid := IsNamingValid("camel")
		return v == namingCamel && valid
	}())
	assert.True(t, func() bool {
		v, valid := IsNamingValid("snake")
		return v == namingSnake && valid
	}())
	assert.True(t, func() bool {
		v, valid := IsNamingValid("")
		return v == namingLower && valid
	}())
	assert.False(t, func() bool {
		v, valid := IsNamingValid("lower ")
		return v == namingLower && valid
	}())
	assert.False(t, func() bool {
		v, valid := IsNamingValid("snake_case")
		return v == namingSnake && valid
	}())

}
