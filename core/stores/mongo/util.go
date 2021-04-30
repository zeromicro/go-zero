package mongo

import "strings"

const mongoAddrSep = ","

// FormatAddr formats mongo hosts to a string.
func FormatAddr(hosts []string) string {
	return strings.Join(hosts, mongoAddrSep)
}
