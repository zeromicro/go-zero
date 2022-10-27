package fs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloseOnExec(t *testing.T) {
	file := os.NewFile(0, os.DevNull)
	assert.NotPanics(t, func() {
		CloseOnExec(file)
	})
}
