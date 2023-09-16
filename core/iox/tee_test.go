package iox

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitTeeReader(t *testing.T) {
	limit := int64(4)
	src := []byte("hello, world")
	dst := make([]byte, len(src))
	rb := bytes.NewBuffer(src)
	wb := new(bytes.Buffer)
	r := LimitTeeReader(rb, wb, limit)
	if n, err := io.ReadFull(r, dst); err != nil || n != len(src) {
		t.Fatalf("ReadFull(r, dst) = %d, %v; want %d, nil", n, err, len(src))
	}
	if !bytes.Equal(dst, src) {
		t.Errorf("bytes read = %q want %q", dst, src)
	}
	if !bytes.Equal(wb.Bytes(), src[:limit]) {
		t.Errorf("bytes written = %q want %q", wb.Bytes(), src)
	}

	n, err := r.Read(dst)
	assert.Equal(t, 0, n)
	assert.Equal(t, io.EOF, err)

	rb = bytes.NewBuffer(src)
	pr, pw := io.Pipe()
	if assert.NoError(t, pr.Close()) {
		r = LimitTeeReader(rb, pw, limit)
		n, err := io.ReadFull(r, dst)
		assert.Equal(t, 0, n)
		assert.Equal(t, io.ErrClosedPipe, err)
	}
}
