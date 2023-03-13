package logx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/color"
)

func TestWithColor(t *testing.T) {
	//old := encoding.Load()
	old := encoding
	SetEncoding(&PlainTextLogEncoder{WithColor: true})
	defer SetEncoding(old)

	output := WithColor("hello", color.BgBlue)
	assert.Equal(t, "hello", output)

	SetEncoding(&JsonLogEncoder{})
	output = WithColor("hello", color.BgBlue)
	assert.Equal(t, "hello", output)
}

func TestWithColorPadding(t *testing.T) {
	//old := encoding.Load()
	old := encoding
	SetEncoding(&PlainTextLogEncoder{WithColor: true})
	defer SetEncoding(old)

	output := WithColorPadding("hello", color.BgBlue)
	assert.Equal(t, " hello ", output)

	SetEncoding(&JsonLogEncoder{})
	output = WithColorPadding("hello", color.BgBlue)
	assert.Equal(t, "hello", output)
}

func TestWithoutColorPadding(t *testing.T) {
	//old := encoding.Load()
	old := encoding
	SetEncoding(&PlainTextLogEncoder{})
	defer SetEncoding(old)

	output := WithColorPadding("hello", color.BgBlue)
	assert.Equal(t, "hello", output)
}
