package breaker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucketAdd(t *testing.T) {
	b := &bucket{}

	// Test succeed
	b.Add(0) // Using 0 for success
	assert.Equal(t, int64(1), b.Sum, "Sum should be incremented")
	assert.Equal(t, int64(1), b.Success, "Success should be incremented")
	assert.Equal(t, int64(0), b.Failure, "Failure should not be incremented")
	assert.Equal(t, int64(0), b.Drop, "Drop should not be incremented")

	// Test failure
	b.Add(fail)
	assert.Equal(t, int64(2), b.Sum, "Sum should be incremented")
	assert.Equal(t, int64(1), b.Failure, "Failure should be incremented")
	assert.Equal(t, int64(0), b.Drop, "Drop should not be incremented")

	// Test drop
	b.Add(drop)
	assert.Equal(t, int64(3), b.Sum, "Sum should be incremented")
	assert.Equal(t, int64(1), b.Drop, "Drop should be incremented")
}

func TestBucketReset(t *testing.T) {
	b := &bucket{
		Sum:     3,
		Success: 1,
		Failure: 1,
		Drop:    1,
	}
	b.Reset()
	assert.Equal(t, int64(0), b.Sum, "Sum should be reset to 0")
	assert.Equal(t, int64(0), b.Success, "Success should be reset to 0")
	assert.Equal(t, int64(0), b.Failure, "Failure should be reset to 0")
	assert.Equal(t, int64(0), b.Drop, "Drop should be reset to 0")
}
