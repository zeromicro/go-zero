package mon

import (
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
)

const mongoAddrSep = ","

// FormatAddr formats mongo hosts to a string.
func FormatAddr(hosts []string) string {
	return strings.Join(hosts, mongoAddrSep)
}

func logDuration(name, method string, startTime time.Duration, err error) {
	duration := timex.Since(startTime)
	if err != nil {
		logx.WithDuration(duration).Infof("mongo(%s) - %s - fail(%s)", name, method, err.Error())
	} else {
		logx.WithDuration(duration).Infof("mongo(%s) - %s - ok", name, method)
	}
}
