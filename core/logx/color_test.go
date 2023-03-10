package logx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/color"
)

func TestWithColor(t *testing.T) {
	//old := encoding.Load()
	old := encoding
	SetEncoding(&PlainTextLogEncoding{})
	defer SetEncoding(old)

	output := WithColor("hello", color.BgBlue)
	assert.Equal(t, "hello", output)

	SetEncoding(&JsonLogEncoding{})
	output = WithColor("hello", color.BgBlue)
	assert.Equal(t, "hello", output)
}

func TestWithColorPadding(t *testing.T) {
	//old := encoding.Load()
	old := encoding
	SetEncoding(&PlainTextLogEncoding{})
	defer SetEncoding(old)

	output := WithColorPadding("hello", color.BgBlue)
	assert.Equal(t, " hello ", output)

	SetEncoding(&JsonLogEncoding{})
	output = WithColorPadding("hello", color.BgBlue)
	assert.Equal(t, "hello", output)
}
