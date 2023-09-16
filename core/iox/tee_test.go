package iox

import (
	"bytes"
	"io"
	"testing"
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
	if n, err := r.Read(dst); n != 0 || err != io.EOF {
		t.Errorf("r.Read at EOF = %d, %v want 0, EOF", n, err)
	}
	rb = bytes.NewBuffer(src)
	pr, pw := io.Pipe()
	pr.Close()
	r = LimitTeeReader(rb, pw, limit)
	if n, err := io.ReadFull(r, dst); n != 0 || err != io.ErrClosedPipe {
		t.Errorf("closed tee: ReadFull(r, dst) = %d, %v; want 0, EPIPE", n, err)
	}
}
