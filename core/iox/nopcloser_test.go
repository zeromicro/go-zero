package iox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNopCloser(t *testing.T) {
	closer := NopCloser(nil)
	assert.NoError(t, closer.Close())
}
