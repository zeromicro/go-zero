package filex

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/fs"
)

func TestSplitLineChunks(t *testing.T) {
	const text = `first line
second line
third line
fourth line
fifth line
sixth line
seventh line
`
	fp, err := fs.TempFileWithText(text)
	assert.Nil(t, err)
	defer func() {
		fp.Close()
		os.Remove(fp.Name())
	}()

	offsets, err := SplitLineChunks(fp.Name(), 3)
	assert.Nil(t, err)
	body := make([]byte, 512)
	for _, offset := range offsets {
		reader := NewRangeReader(fp, offset.Start, offset.Stop)
		n, err := reader.Read(body)
		assert.Nil(t, err)
		assert.Equal(t, uint8('\n'), body[n-1])
	}
}

func TestSplitLineChunksNoFile(t *testing.T) {
	_, err := SplitLineChunks("nosuchfile", 2)
	assert.NotNil(t, err)
}

func TestSplitLineChunksFull(t *testing.T) {
	const text = `first line
second line
third line
fourth line
fifth line
sixth line
`
	fp, err := fs.TempFileWithText(text)
	assert.Nil(t, err)
	defer func() {
		fp.Close()
		os.Remove(fp.Name())
	}()

	offsets, err := SplitLineChunks(fp.Name(), 1)
	assert.Nil(t, err)
	body := make([]byte, 512)
	for _, offset := range offsets {
		reader := NewRangeReader(fp, offset.Start, offset.Stop)
		n, err := reader.Read(body)
		assert.Nil(t, err)
		assert.Equal(t, []byte(text), body[:n])
	}
}
