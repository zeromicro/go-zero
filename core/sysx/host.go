package sysx

import (
	"os"

	"zero/core/lang"
)

var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	lang.Must(err)
}

func Hostname() string {
	return hostname
}
