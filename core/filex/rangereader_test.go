package filex

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/fs"
)

func TestRangeReader(t *testing.T) {
	const text = `hello
world`
	file, err := fs.TempFileWithText(text)
	assert.Nil(t, err)
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	reader := NewRangeReader(file, 5, 8)
	buf := make([]byte, 10)
	n, err := reader.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, `
wo`, string(buf[:n]))

	n, err = reader.Read(buf)
	assert.Zero(t, n)
	assert.ErrorIs(t, err, io.EOF)
}

func TestRangeReader_OutOfRange(t *testing.T) {
	const text = `hello
world`
	file, err := fs.TempFileWithText(text)
	assert.Nil(t, err)
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	reader := NewRangeReader(file, 50, 8)
	buf := make([]byte, 10)
	_, err = reader.Read(buf)
	assert.NotNil(t, err)
}

func TestRangeReader_AtFileEOF(t *testing.T) {
	const text = "hello\nworld"
	file, err := fs.TempFileWithText(text)
	if !assert.NoError(t, err) {
		return
	}
	t.Cleanup(func() {
		file.Close()
		os.Remove(file.Name())
	})

	for _, test := range []struct {
		name  string
		start int64
		want  string
		err   error
	}{
		{name: "range ending at file EOF", start: 5, want: "\nworld"},
		{name: "empty range at file EOF", start: int64(len(text)), err: io.EOF},
	} {
		t.Run(test.name, func(t *testing.T) {
			reader := NewRangeReader(file, test.start, int64(len(text)))
			buf := make([]byte, len(text)+1)
			n, err := reader.Read(buf)
			assert.Equal(t, len(test.want), n)
			assert.Equal(t, test.want, string(buf[:n]))
			assert.ErrorIs(t, err, test.err)

			for i := 0; i < 2; i++ {
				n, err = reader.Read(buf)
				assert.Zero(t, n)
				assert.ErrorIs(t, err, io.EOF)
			}
		})
	}
}
