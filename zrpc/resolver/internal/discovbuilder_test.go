package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscovBuilder_Scheme(t *testing.T) {
	var b discovBuilder
	assert.Equal(t, DiscovScheme, b.Scheme())
}
