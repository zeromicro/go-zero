package logx

import (
	"sync/atomic"

	"github.com/zeromicro/go-zero/core/color"
)

// WithColor is a helper function to add color to a string, only in plain encoding.
func WithColor(text, colour string) string {
	if atomic.LoadUint32(&encoding) == plainEncodingType {
		return color.WithColor(text, colour)
	}

	return text
}

// WithColorPadding is a helper function to add color to a string with leading and trailing spaces,
// only in plain encoding.
func WithColorPadding(text, colour string) string {
	if atomic.LoadUint32(&encoding) == plainEncodingType {
		return color.WithColorPadding(text, colour)
	}

	return text
}
