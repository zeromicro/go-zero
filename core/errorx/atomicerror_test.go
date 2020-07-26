package errorx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errDummy = errors.New("hello")

func TestAtomicError(t *testing.T) {
	var err AtomicError
	err.Set(errDummy)
	assert.Equal(t, errDummy, err.Load())
}

func TestAtomicErrorNil(t *testing.T) {
	var err AtomicError
	assert.Nil(t, err.Load())
}
