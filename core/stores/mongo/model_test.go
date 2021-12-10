package mongo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithTimeout(t *testing.T) {
	o := defaultOptions()
	WithTimeout(time.Second)(o)
	assert.Equal(t, time.Second, o.timeout)
}
