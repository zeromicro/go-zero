package color

const (
	Reset = "\033[0m"

	Black   = "\033[97;30m"
	Red     = "\033[97;31m"
	Green   = "\033[97;32m"
	Yellow  = "\033[97;33m"
	Blue    = "\033[97;34m"
	Magenta = "\033[97;35m"
	Cyan    = "\033[97;36m"
	White   = "\033[97;37m"

	BgBlack   = "\033[97;40m"
	BgRed     = "\033[97;41m"
	BgGreen   = "\033[97;42m"
	BgYellow  = "\033[97;43m"
	BgBlue    = "\033[97;44m"
	BgMagenta = "\033[97;45m"
	BgCyan    = "\033[97;46m"
	BgWhite   = "\033[97;47m"
)

func WithColor(text, color string) string {
	return color + text + Reset
}

func WithColorPadding(text, color string) string {
	return color + " " + text + " " + Reset
}
