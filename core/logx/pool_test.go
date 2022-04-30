package logx

import "testing"

func TestPool_getBuffer(t *testing.T) {
	for i := 0; i < 100; i++ {
		b := getBuffer()
		b.WriteString("hello")
		putBuffer(b)
	}

	b := getBuffer()
	if b == nil {
		t.Error("getBuffer() failed")
	}
	if b.Len() != 0 {
		t.Error("getBuffer() failed")
	}
	putBuffer(b)
}
