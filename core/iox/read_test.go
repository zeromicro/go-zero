package iox

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/fs"
	"github.com/zeromicro/go-zero/core/stringx"
)

func TestReadText(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			input:  `a`,
			expect: `a`,
		}, {
			input: `a
`,
			expect: `a`,
		}, {
			input: `a
b`,
			expect: `a
b`,
		}, {
			input: `a
b
`,
			expect: `a
b`,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			tmpFile, err := fs.TempFilenameWithText(test.input)
			assert.Nil(t, err)
			defer os.Remove(tmpFile)

			content, err := ReadText(tmpFile)
			assert.Nil(t, err)
			assert.Equal(t, test.expect, content)
		})
	}
}

func TestReadTextError(t *testing.T) {
	_, err := ReadText("not-exist")
	assert.NotNil(t, err)
}

func TestReadTextLines(t *testing.T) {
	text := `1

    2

    #a
    3`

	tmpFile, err := fs.TempFilenameWithText(text)
	assert.Nil(t, err)
	defer os.Remove(tmpFile)

	tests := []struct {
		options     []TextReadOption
		expectLines int
	}{
		{
			nil,
			6,
		}, {
			[]TextReadOption{KeepSpace(), OmitWithPrefix("#")},
			6,
		}, {
			[]TextReadOption{WithoutBlank()},
			4,
		}, {
			[]TextReadOption{OmitWithPrefix("#")},
			5,
		}, {
			[]TextReadOption{WithoutBlank(), OmitWithPrefix("#")},
			3,
		},
	}

	for _, test := range tests {
		t.Run(stringx.Rand(), func(t *testing.T) {
			lines, err := ReadTextLines(tmpFile, test.options...)
			assert.Nil(t, err)
			assert.Equal(t, test.expectLines, len(lines))
		})
	}
}

func TestReadTextLinesError(t *testing.T) {
	_, err := ReadTextLines("not-exist")
	assert.NotNil(t, err)
}

func TestDupReadCloser(t *testing.T) {
	input := "hello"
	reader := io.NopCloser(bytes.NewBufferString(input))
	r1, r2 := DupReadCloser(reader)
	verify := func(r io.Reader) {
		output, err := io.ReadAll(r)
		assert.Nil(t, err)
		assert.Equal(t, input, string(output))
	}

	verify(r1)
	verify(r2)
}

func TestLimitDupReadCloser(t *testing.T) {
	input := "hello world"
	limitBytes := int64(4)
	reader := io.NopCloser(bytes.NewBufferString(input))
	r1, r2 := LimitDupReadCloser(reader, limitBytes)
	verify := func(r io.Reader) {
		output, err := io.ReadAll(r)
		assert.Nil(t, err)
		assert.Equal(t, input, string(output))
	}
	verifyLimit := func(r io.Reader, limit int64) {
		output, err := io.ReadAll(r)
		if limit < int64(len(input)) {
			input = input[:limit]
		}
		assert.Nil(t, err)
		assert.Equal(t, input, string(output))
	}

	verify(r1)
	verifyLimit(r2, limitBytes)
}

func TestReadBytes(t *testing.T) {
	reader := io.NopCloser(bytes.NewBufferString("helloworld"))
	buf := make([]byte, 5)
	err := ReadBytes(reader, buf)
	assert.Nil(t, err)
	assert.Equal(t, "hello", string(buf))
}

func TestReadBytesNotEnough(t *testing.T) {
	reader := io.NopCloser(bytes.NewBufferString("hell"))
	buf := make([]byte, 5)
	err := ReadBytes(reader, buf)
	assert.Equal(t, io.EOF, err)
}

func TestReadBytesChunks(t *testing.T) {
	buf := make([]byte, 5)
	reader, writer := io.Pipe()

	go func() {
		for i := 0; i < 10; i++ {
			writer.Write([]byte{'a'})
			time.Sleep(10 * time.Millisecond)
		}
	}()

	err := ReadBytes(reader, buf)
	assert.Nil(t, err)
	assert.Equal(t, "aaaaa", string(buf))
}
