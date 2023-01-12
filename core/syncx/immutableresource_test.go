package syncx

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestImmutableResource(t *testing.T) {
	var count int
	r := NewImmutableResource(func() (any, error) {
		count++
		return "hello", nil
	})

	res, err := r.Get()
	assert.Equal(t, "hello", res)
	assert.Equal(t, 1, count)
	assert.Nil(t, err)

	// again
	res, err = r.Get()
	assert.Equal(t, "hello", res)
	assert.Equal(t, 1, count)
	assert.Nil(t, err)
}

func TestImmutableResourceError(t *testing.T) {
	var count int
	r := NewImmutableResource(func() (any, error) {
		count++
		return nil, errors.New("any")
	})

	res, err := r.Get()
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, "any", err.Error())
	assert.Equal(t, 1, count)

	// again
	res, err = r.Get()
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, "any", err.Error())
	assert.Equal(t, 1, count)

	r.refreshInterval = 0
	time.Sleep(time.Millisecond)
	res, err = r.Get()
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, "any", err.Error())
	assert.Equal(t, 2, count)
}

func TestImmutableResourceErrorRefreshAlways(t *testing.T) {
	var count int
	r := NewImmutableResource(func() (any, error) {
		count++
		return nil, errors.New("any")
	}, WithRefreshIntervalOnFailure(0))

	res, err := r.Get()
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, "any", err.Error())
	assert.Equal(t, 1, count)

	// again
	res, err = r.Get()
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, "any", err.Error())
	assert.Equal(t, 2, count)
}
