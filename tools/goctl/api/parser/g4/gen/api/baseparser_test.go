package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	assert.False(t, matchRegex("v1ddd", versionRegex))
}
