package errorx

import (
	"errors"
	"fmt"
	"sync"
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
	assert.NotNil(t, batch.Err())
	assert.Equal(t, err1, batch.Err().Error())
	assert.True(t, batch.NotNil())
}

func TestBatchErrorWithErrors(t *testing.T) {
	var batch BatchError
	batch.Add(errors.New(err1))
	batch.Add(errors.New(err2))
	assert.NotNil(t, batch.Err())
	assert.Equal(t, fmt.Sprintf("%s\n%s", err1, err2), batch.Err().Error())
	assert.True(t, batch.NotNil())
}

func TestBatchErrorConcurrentAdd(t *testing.T) {
	const count = 10000
	var batch BatchError
	var wg sync.WaitGroup

	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			batch.Add(errors.New(err1))
		}()
	}
	wg.Wait()

	assert.NotNil(t, batch.Err())
	assert.Equal(t, count, len(batch.errs))
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

func TestBatchError_Add(t *testing.T) {
	var be BatchError

	// Test adding nil errors
	be.Add(nil, nil)
	assert.False(t, be.NotNil(), "Expected BatchError to be empty after adding nil errors")

	// Test adding non-nil errors
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	be.Add(err1, err2)
	assert.True(t, be.NotNil(), "Expected BatchError to be non-empty after adding errors")

	// Test adding a mix of nil and non-nil errors
	err3 := errors.New("error 3")
	be.Add(nil, err3, nil)
	assert.True(t, be.NotNil(), "Expected BatchError to be non-empty after adding a mix of nil and non-nil errors")
}

func TestBatchError_Err(t *testing.T) {
	var be BatchError

	// Test Err() on empty BatchError
	assert.Nil(t, be.Err(), "Expected nil error for empty BatchError")

	// Test Err() with multiple errors
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	be.Add(err1, err2)

	combinedErr := be.Err()
	assert.NotNil(t, combinedErr, "Expected nil error for BatchError with multiple errors")

	// Check if the combined error contains both error messages
	errString := combinedErr.Error()
	assert.Truef(t, errors.Is(combinedErr, err1), "Combined error doesn't contain first error: %s", errString)
	assert.Truef(t, errors.Is(combinedErr, err2), "Combined error doesn't contain second error: %s", errString)
}

func TestBatchError_NotNil(t *testing.T) {
	var be BatchError

	// Test NotNil() on empty BatchError
	assert.Nil(t, be.Err(), "Expected nil error for empty BatchError")

	// Test NotNil() after adding an error
	be.Add(errors.New("test error"))
	assert.NotNil(t, be.Err(), "Expected non-nil error after adding an error")
}
