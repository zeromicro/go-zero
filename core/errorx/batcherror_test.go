package errorx

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	err1 = "first error"
	err2 = "second error"
)

func TestBatchErrorNil(t *testing.T) {
	var batch BatchError
	assert.Nil(t, batch.Err())
	assert.False(t, batch.NotNil())
	batch.Add(nil)
	assert.Nil(t, batch.Err())
	assert.False(t, batch.NotNil())
}

func TestBatchErrorNilFromFunc(t *testing.T) {
	err := func() error {
		var be BatchError
		return be.Err()
	}()
	assert.True(t, err == nil)
}

func TestBatchErrorOneError(t *testing.T) {
	var batch BatchError
	batch.Add(errors.New(err1))
	assert.NotNil(t, batch)
	assert.Equal(t, err1, batch.Err().Error())
	assert.True(t, batch.NotNil())
}

func TestBatchErrorWithErrors(t *testing.T) {
	var batch BatchError
	batch.Add(errors.New(err1))
	batch.Add(errors.New(err2))
	assert.NotNil(t, batch)
	assert.Equal(t, fmt.Sprintf("%s\n%s", err1, err2), batch.Err().Error())
	assert.True(t, batch.NotNil())
}

func TestBatchError_Unwrap(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var be BatchError
		assert.Nil(t, be.Err())
		assert.True(t, errors.Is(be.Err(), nil))
	})

	t.Run("one error", func(t *testing.T) {
		var errFoo = errors.New("foo")
		var errBar = errors.New("bar")
		var be BatchError
		be.Add(errFoo)
		assert.True(t, errors.Is(be.Err(), errFoo))
		assert.False(t, errors.Is(be.Err(), errBar))
	})

	t.Run("two errors", func(t *testing.T) {
		var errFoo = errors.New("foo")
		var errBar = errors.New("bar")
		var errBaz = errors.New("baz")
		var be BatchError
		be.Add(errFoo)
		be.Add(errBar)
		assert.True(t, errors.Is(be.Err(), errFoo))
		assert.True(t, errors.Is(be.Err(), errBar))
		assert.False(t, errors.Is(be.Err(), errBaz))
	})
}
