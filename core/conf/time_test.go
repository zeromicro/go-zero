package conf

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckedDuration(t *testing.T) {
	assert.Equal(t, time.Second, CheckedDuration(1000))
	assert.Equal(t, 2*time.Second, CheckedDuration(2000))
	assert.Equal(t, 2*time.Second, CheckedDuration(time.Second*2))
}
