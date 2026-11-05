package filex

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/fs"
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

func TestSplitLineChunksLastLineLongerThanPreferSize(t *testing.T) {
	// preferSize ≈ size/chunks+1; with chunks=2 and a short first line
	// plus a long last line, the ideal cut lands inside the last line.
	content := "a\n" + strings.Repeat("b", 3000)
	fp, err := fs.TempFileWithText(content)
	assert.Nil(t, err)
	defer func() {
		fp.Close()
		os.Remove(fp.Name())
	}()

	ranges, err := SplitLineChunks(fp.Name(), 2)
	assert.Nil(t, err)
	assert.NotEmpty(t, ranges)

	var covered int64
	for _, r := range ranges {
		assert.GreaterOrEqual(t, r.Stop, r.Start)
		covered += r.Stop - r.Start
	}
	assert.Equal(t, int64(len(content)), covered)
	assert.Equal(t, int64(0), ranges[0].Start)
	assert.Equal(t, int64(len(content)), ranges[len(ranges)-1].Stop)
}
