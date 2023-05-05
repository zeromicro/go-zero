package cmd

import (
	"fmt"

	"github.com/gookit/color"
)

var colorRender = []func(v any) string{
	func(v any) string {
		return color.LightRed.Render(v)
	},
	func(v any) string {
		return color.LightGreen.Render(v)
	},
	func(v any) string {
		return color.LightYellow.Render(v)
	},
	func(v any) string {
		return color.LightBlue.Render(v)
	},
	func(v any) string {
		return color.LightMagenta.Render(v)
	},
	func(v any) string {
		return color.LightCyan.Render(v)
	},
}

func blue(s string) string {
	return color.LightBlue.Render(s)
}

func green(s string) string {
	return color.LightGreen.Render(s)
}

func rainbow(s string) string {
	s0 := s[0]
	return colorRender[int(s0)%(len(colorRender)-1)](s)
}

// rpadx adds padding to the right of a string.
func rpadx(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds", padding)
	return rainbow(fmt.Sprintf(template, s))
}
