package sysx

import (
	"os"

	"github.com/tal-tech/go-zero/core/lang"
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
