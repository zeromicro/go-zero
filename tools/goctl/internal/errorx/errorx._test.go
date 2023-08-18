package errorx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	err := errors.New("foo")
	err = Wrap(err)
	var goctlError *GoctlError
	ok := errors.As(err, &goctlError)
	assert.True(t, ok)
}
