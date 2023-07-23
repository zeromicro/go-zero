package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	assert.NotPanics(t, func() {
		Register()
	})
}
