package proc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessName(t *testing.T) {
	assert.True(t, len(ProcessName()) > 0)
}

func TestPid(t *testing.T) {
	assert.True(t, Pid() > 0)
}
