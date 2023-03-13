package logx

import (
	"github.com/zeromicro/go-zero/core/color"
)

// LogColor parse log text color tag, see `github.com/fatih/color`
type LogColor interface {
	Color() bool
}

// WithColor is a helper function to add color to a string, only in plain encoding.
func WithColor(text string, colour color.Color) string {
	if c, ok := encoding.(LogColor); ok && c.Color() {
		return color.WithColor(text, colour)
	}

	return text
}

// WithColorPadding is a helper function to add color to a string with leading and trailing spaces,
// only in plain encoding.
func WithColorPadding(text string, colour color.Color) string {
	if c, ok := encoding.(LogColor); ok && c.Color() {
		return color.WithColorPadding(text, colour)
	}

	return text
}
