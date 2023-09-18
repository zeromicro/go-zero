package iox

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountLines(t *testing.T) {
	const val = `1
2
3
4`
	file, err := os.CreateTemp(os.TempDir(), "test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	file.WriteString(val)
	file.Close()
	lines, err := CountLines(file.Name())
	assert.Nil(t, err)
	assert.Equal(t, 4, lines)
}

func TestCountLinesError(t *testing.T) {
	_, err := CountLines("not-exist")
	assert.NotNil(t, err)
}
