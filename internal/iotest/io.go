package iotest

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func RunTest(t *testing.T, fn func(), validate func(string)) {
	defer func(orig *os.File) {
		os.Stdout = orig
	}(os.Stdout)

	r, w, err := os.Pipe()
	assert.NoError(t, err)

	os.Stdout = w
	fn()
	assert.NoError(t, w.Close())

	out, _ := io.ReadAll(r)
	validate(string(out))
}
