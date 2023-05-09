package logtest

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/zeromicro/go-zero/core/logx"
)

type Buffer struct {
	buf *bytes.Buffer
	t   *testing.T
}

func Discard(t *testing.T) {
	prev := logx.Reset()
	logx.SetWriter(logx.NewWriter(io.Discard))

	t.Cleanup(func() {
		logx.SetWriter(prev)
	})
}

func NewCollector(t *testing.T) *Buffer {
	var buf bytes.Buffer
	writer := logx.NewWriter(&buf)
	prev := logx.Reset()
	logx.SetWriter(writer)

	t.Cleanup(func() {
		logx.SetWriter(prev)
	})

	return &Buffer{
		buf: &buf,
		t:   t,
	}
}

func (b *Buffer) Bytes() []byte {
	return b.buf.Bytes()
}

func (b *Buffer) Content() string {
	var m map[string]interface{}
	if err := json.Unmarshal(b.buf.Bytes(), &m); err != nil {
		return ""
	}

	content, ok := m["content"]
	if !ok {
		return ""
	}

	switch val := content.(type) {
	case string:
		return val
	default:
		// err is impossible to be not nil, unmarshaled from b.buf.Bytes()
		bs, _ := json.Marshal(content)
		return string(bs)
	}
}

func (b *Buffer) Reset() {
	b.buf.Reset()
}

func (b *Buffer) String() string {
	return b.buf.String()
}

func PanicOnFatal(t *testing.T) {
	ok := logx.ExitOnFatal.CompareAndSwap(true, false)
	if !ok {
		return
	}

	t.Cleanup(func() {
		logx.ExitOnFatal.CompareAndSwap(false, true)
	})
}
