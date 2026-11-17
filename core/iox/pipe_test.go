package iox

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirectInOut(t *testing.T) {
	originalStdin, originalStdout := os.Stdin, os.Stdout
	restore, err := RedirectInOut()
	require.NoError(t, err)

	stdin, stdout := os.Stdin, os.Stdout
	t.Cleanup(func() {
		restore()
		// Also clean up when testing an implementation that leaks the pipe.
		_ = stdin.Close()
		_ = stdout.Close()
	})

	restore()
	assert.Same(t, originalStdin, os.Stdin)
	assert.Same(t, originalStdout, os.Stdout)
	_, err = stdin.Stat()
	assert.ErrorIs(t, err, os.ErrClosed)
	_, err = stdout.Stat()
	assert.ErrorIs(t, err, os.ErrClosed)
}
