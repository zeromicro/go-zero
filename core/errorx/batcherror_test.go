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
