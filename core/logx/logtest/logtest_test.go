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
}

func TestPanicOnFatal(t *testing.T) {
	const input = "hello"
	Discard(t)
	logx.Info(input)

	PanicOnFatal(t)
	assert.Panics(t, func() {
		logx.Must(errors.New("foo"))
	})
}
