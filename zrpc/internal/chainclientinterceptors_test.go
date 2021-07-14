package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithStreamClientInterceptors(t *testing.T) {
	opts := WithStreamClientInterceptors()
	assert.NotNil(t, opts)
}

func TestWithUnaryClientInterceptors(t *testing.T) {
	opts := WithUnaryClientInterceptors()
	assert.NotNil(t, opts)
}
