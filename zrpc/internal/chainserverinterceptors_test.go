package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithStreamServerInterceptors(t *testing.T) {
	opts := WithStreamServerInterceptors()
	assert.NotNil(t, opts)
}

func TestWithUnaryServerInterceptors(t *testing.T) {
	opts := WithUnaryServerInterceptors()
	assert.NotNil(t, opts)
}
