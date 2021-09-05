package errorx

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWarp(t *testing.T) {
	err := errors.New("foo")
	err = Wrap(err)
	_, ok := err.(*GoctlError)
	assert.True(t, ok)
}
