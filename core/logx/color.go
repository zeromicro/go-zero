package logx

import (
	"github.com/zeromicro/go-zero/core/color"
)

// WithColor is a helper function to add color to a string, only in plain encoding.
func WithColor(text string, colour color.Color) string {
	if _, ok := encoding.(*PlainTextLogEncoder); ok {
		return color.WithColor(text, colour)
	}

	return text
}

// WithColorPadding is a helper function to add color to a string with leading and trailing spaces,
// only in plain encoding.
func WithColorPadding(text string, colour color.Color) string {
	if _, ok := encoding.(*PlainTextLogEncoder); ok {
		return color.WithColorPadding(text, colour)
	}

	return text
}
