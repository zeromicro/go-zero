package color

import "github.com/fatih/color"

const (
	// NoColor is no color for both foreground and background.
	NoColor Color = iota
	// FgBlack is the foreground color black.
	FgBlack
	// FgRed is the foreground color red.
	FgRed
	// FgGreen is the foreground color green.
	FgGreen
	// FgYellow is the foreground color yellow.
	FgYellow
	// FgBlue is the foreground color blue.
	FgBlue
	// FgMagenta is the foreground color magenta.
	FgMagenta
	// FgCyan is the foreground color cyan.
	FgCyan
	// FgWhite is the foreground color white.
	FgWhite

	// BgBlack is the background color black.
	BgBlack
	// BgRed is the background color red.
	BgRed
	// BgGreen is the background color green.
	BgGreen
	// BgYellow is the background color yellow.
	BgYellow
	// BgBlue is the background color blue.
	BgBlue
	// BgMagenta is the background color magenta.
	BgMagenta
	// BgCyan is the background color cyan.
	BgCyan
	// BgWhite is the background color white.
	BgWhite
)

var colors = map[Color][]color.Attribute{
	FgBlack:   {color.FgBlack, color.Bold},
	FgRed:     {color.FgRed, color.Bold},
	FgGreen:   {color.FgGreen, color.Bold},
	FgYellow:  {color.FgYellow, color.Bold},
	FgBlue:    {color.FgBlue, color.Bold},
	FgMagenta: {color.FgMagenta, color.Bold},
	FgCyan:    {color.FgCyan, color.Bold},
	FgWhite:   {color.FgWhite, color.Bold},
	BgBlack:   {color.BgBlack, color.FgHiWhite, color.Bold},
	BgRed:     {color.BgRed, color.FgHiWhite, color.Bold},
	BgGreen:   {color.BgGreen, color.FgHiWhite, color.Bold},
	BgYellow:  {color.BgHiYellow, color.FgHiBlack, color.Bold},
	BgBlue:    {color.BgBlue, color.FgHiWhite, color.Bold},
	BgMagenta: {color.BgMagenta, color.FgHiWhite, color.Bold},
	BgCyan:    {color.BgCyan, color.FgHiWhite, color.Bold},
	BgWhite:   {color.BgHiWhite, color.FgHiBlack, color.Bold},
}

type Color uint32

// WithColor returns a string with the given color applied.
func WithColor(text string, colour Color) string {
	c := color.New(colors[colour]...)
	return c.Sprint(text)
}

// WithColorPadding returns a string with the given color applied with leading and trailing spaces.
func WithColorPadding(text string, colour Color) string {
	return WithColor(" "+text+" ", colour)
}
