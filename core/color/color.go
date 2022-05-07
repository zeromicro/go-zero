package color

const (
	// Reset restores the default terminal foreground and background colors.
	Reset = "\033[0m"

	// Black is the color code for foreground black.
	Black = "\033[97;30m"
	// Red is the color code for foreground red.
	Red = "\033[97;31m"
	// Green is the color code for foreground green.
	Green = "\033[97;32m"
	// Yellow is the color code for foreground yellow.
	Yellow = "\033[97;33m"
	// Blue is the color code for foreground blue.
	Blue = "\033[97;34m"
	// Magenta is the color code for foreground magenta.
	Magenta = "\033[97;35m"
	// Cyan is the color code for foreground cyan.
	Cyan = "\033[97;36m"
	// White is the color code for foreground white.
	White = "\033[97;37m"

	// BgBlack is the color code for background black.
	BgBlack = "\033[97;40m"
	// BgRed is the color code for background red.
	BgRed = "\033[97;41m"
	// BgGreen is the color code for background green.
	BgGreen = "\033[97;42m"
	// BgYellow is the color code for background yellow.
	BgYellow = "\033[97;43m"
	// BgBlue is the color code for background blue.
	BgBlue = "\033[97;44m"
	// BgMagenta is the color code for background magenta.
	BgMagenta = "\033[97;45m"
	// BgCyan is the color code for background cyan.
	BgCyan = "\033[97;46m"
	// BgWhite is the color code for background white.
	BgWhite = "\033[97;47m"
)

// WithColor returns a string with the given color applied.
func WithColor(text, color string) string {
	return color + text + Reset
}

// WithColorPadding returns a string with the given color applied with leading and trailing spaces.
func WithColorPadding(text, color string) string {
	return color + " " + text + " " + Reset
}
