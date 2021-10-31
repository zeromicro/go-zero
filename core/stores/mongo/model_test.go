package mongo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithSlowThreshold(t *testing.T) {
	o := defaultOptions()
	WithSlowThreshold(time.Second)(o)
	assert.Equal(t, time.Second, o.slowThreshold)
}

func TestWithTimeout(t *testing.T) {
	o := defaultOptions()
	WithTimeout(time.Second)(o)
	assert.Equal(t, time.Second, o.timeout)
}
