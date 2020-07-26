package syncx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnceGuard(t *testing.T) {
	var guard OnceGuard

	assert.False(t, guard.Taken())
	assert.True(t, guard.Take())
	assert.True(t, guard.Taken())
	assert.False(t, guard.Take())
	assert.True(t, guard.Taken())
}
