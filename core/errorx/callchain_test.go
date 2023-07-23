package errorx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	errDummy := errors.New("dummy")
	assert.Nil(t, Chain(func() error {
		return nil
	}, func() error {
		return nil
	}))
	assert.Equal(t, errDummy, Chain(func() error {
		return errDummy
	}, func() error {
		return nil
	}))
	assert.Equal(t, errDummy, Chain(func() error {
		return nil
	}, func() error {
		return errDummy
	}))
}
