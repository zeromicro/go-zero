package mongo

import "strings"

const mongoAddrSep = ","

func FormatAddr(hosts []string) string {
	return strings.Join(hosts, mongoAddrSep)
}
