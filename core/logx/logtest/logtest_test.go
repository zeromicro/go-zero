package logtest

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestCollector(t *testing.T) {
	const input = "hello"
	c := NewCollector(t)
	logx.Info(input)
	assert.Equal(t, input, c.Content())
	assert.Contains(t, c.String(), input)
	c.Reset()
	assert.Empty(t, c.Bytes())
}

func TestPanicOnFatal(t *testing.T) {
	const input = "hello"
	Discard(t)
	logx.Info(input)

	PanicOnFatal(t)
	PanicOnFatal(t)
	assert.Panics(t, func() {
		logx.Must(errors.New("foo"))
	})
}

func TestCollectorContent(t *testing.T) {
	const input = "hello"
	c := NewCollector(t)
	c.buf.WriteString(input)
	assert.Empty(t, c.Content())
	c.Reset()
	c.buf.WriteString(`{}`)
	assert.Empty(t, c.Content())
	c.Reset()
	c.buf.WriteString(`{"content":1}`)
	assert.Equal(t, "1", c.Content())
}
