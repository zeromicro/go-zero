package iox

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedirectInOut(t *testing.T) {
	restore, err := RedirectInOut()
	assert.Nil(t, err)

	stdin := os.Stdin
	stdout := os.Stdout
	defer func() {
		restore()
		_ = stdin.Close()
		_ = stdout.Close()
	}()

	restore()
	_, err = stdin.Stat()
	assert.ErrorIs(t, err, os.ErrClosed)
	_, err = stdout.Stat()
	assert.ErrorIs(t, err, os.ErrClosed)
}
