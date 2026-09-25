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

func TestSplitLineChunksLongLastLine(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		chunks int
		expect [][2]int64
	}{
		{
			name:   "cut point inside the only long line",
			text:   "a\n" + strings.Repeat("b", 3000),
			chunks: 2,
			expect: [][2]int64{{0, 3002}},
		},
		{
			name:   "cut point inside the last line after a full chunk",
			text:   strings.Repeat(strings.Repeat("a", 99)+"\n", 30) + strings.Repeat("b", 3000),
			chunks: 4,
			expect: [][2]int64{{0, 1600}, {1600, 6000}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fp, err := fs.TempFileWithText(test.text)
			assert.Nil(t, err)
			defer func() {
				fp.Close()
				os.Remove(fp.Name())
			}()

			offsets, err := SplitLineChunks(fp.Name(), test.chunks)
			assert.Nil(t, err)
			var actual [][2]int64
			for _, offset := range offsets {
				actual = append(actual, [2]int64{offset.Start, offset.Stop})
			}
			assert.Equal(t, test.expect, actual)
		})
	}
}
