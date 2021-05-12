package util

import "strings"

// TrimStringSlice returns a copy slice without empty string item
func TrimStringSlice(list []string) []string {
	var out []string
	for _, item := range list {
		if len(item) == 0 {
			continue
		}
		out = append(out, item)
	}
	return out
}

func TrimNewLine(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}
