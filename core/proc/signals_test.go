//go:build linux || darwin || freebsd

package proc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDone(t *testing.T) {
	select {
	case <-Done():
		assert.Fail(t, "should run")
	default:
	}
	assert.NotNil(t, Done())
}
