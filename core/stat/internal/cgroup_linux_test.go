package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunningInUserNS(t *testing.T) {
	// should be false in docker
	assert.False(t, runningInUserNS())
}
