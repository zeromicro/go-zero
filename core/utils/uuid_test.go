package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUuid(t *testing.T) {
	assert.Equal(t, 36, len(NewUuid()))
}
