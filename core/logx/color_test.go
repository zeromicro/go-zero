package logx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/color"
)

func TestWithColor(t *testing.T) {
	//old := encoding.Load()
	old := encoding
	setEncoding(&PlainTextLogEncoder{WithColor: true})
	defer setEncoding(old)

	output := WithColor("hello", color.BgBlue)
	assert.Equal(t, "hello", output)

	setEncoding(&JsonLogEncoder{})
	output = WithColor("hello", color.BgBlue)
	assert.Equal(t, "hello", output)
}

func TestWithColorPadding(t *testing.T) {
	//old := encoding.Load()
	old := encoding
	setEncoding(&PlainTextLogEncoder{WithColor: true})
	defer setEncoding(old)

	output := WithColorPadding("hello", color.BgBlue)
	assert.Equal(t, " hello ", output)

	setEncoding(&JsonLogEncoder{})
	output = WithColorPadding("hello", color.BgBlue)
	assert.Equal(t, "hello", output)
}

func TestWithoutColorPadding(t *testing.T) {
	//old := encoding.Load()
	old := encoding
	setEncoding(&PlainTextLogEncoder{})
	defer setEncoding(old)

	output := WithColorPadding("hello", color.BgBlue)
	assert.Equal(t, "hello", output)
}
