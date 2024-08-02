package util

import "strings"

// TrimNewLine trims \r and \n chars.
func TrimNewLine(s string) string {
	return strings.NewReplacer("\r", "", "\n", "").Replace(s)
}
