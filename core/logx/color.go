package logx

import (
	"sync/atomic"

	"github.com/zeromicro/go-zero/core/color"
)

func WithColor(text, colour string) string {
	if atomic.LoadUint32(&encoding) == plainEncodingType {
		return color.WithColor(text, colour)
	}

	return text
}

func WithColorPadding(text, colour string) string {
	if atomic.LoadUint32(&encoding) == plainEncodingType {
		return color.WithColorPadding(text, colour)
	}

	return text
}
