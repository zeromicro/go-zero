package syncx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAtomicDuration(t *testing.T) {
	d := ForAtomicDuration(time.Duration(100))
	assert.Equal(t, time.Duration(100), d.Load())
	d.Set(time.Duration(200))
	assert.Equal(t, time.Duration(200), d.Load())
	assert.True(t, d.CompareAndSwap(time.Duration(200), time.Duration(300)))
	assert.Equal(t, time.Duration(300), d.Load())
	assert.False(t, d.CompareAndSwap(time.Duration(200), time.Duration(400)))
	assert.Equal(t, time.Duration(300), d.Load())
}
