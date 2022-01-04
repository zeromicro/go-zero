package iox

import (
	"bytes"
	"io"
	"io/ioutil"
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
			tmpfile, err := fs.TempFilenameWithText(test.input)
			assert.Nil(t, err)
			defer os.Remove(tmpfile)

			content, err := ReadText(tmpfile)
			assert.Nil(t, err)
			assert.Equal(t, test.expect, content)
		})
	}
}

func TestReadTextLines(t *testing.T) {
	text := `1

    2

    #a
    3`

	tmpfile, err := fs.TempFilenameWithText(text)
	assert.Nil(t, err)
	defer os.Remove(tmpfile)

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
			lines, err := ReadTextLines(tmpfile, test.options...)
			assert.Nil(t, err)
			assert.Equal(t, test.expectLines, len(lines))
		})
	}
}

func TestDupReadCloser(t *testing.T) {
	input := "hello"
	reader := ioutil.NopCloser(bytes.NewBufferString(input))
	r1, r2 := DupReadCloser(reader)
	verify := func(r io.Reader) {
		output, err := ioutil.ReadAll(r)
		assert.Nil(t, err)
		assert.Equal(t, input, string(output))
	}

	verify(r1)
	verify(r2)
}

func TestReadBytes(t *testing.T) {
	reader := ioutil.NopCloser(bytes.NewBufferString("helloworld"))
	buf := make([]byte, 5)
	err := ReadBytes(reader, buf)
	assert.Nil(t, err)
	assert.Equal(t, "hello", string(buf))
}

func TestReadBytesNotEnough(t *testing.T) {
	reader := ioutil.NopCloser(bytes.NewBufferString("hell"))
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
