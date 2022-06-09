package mon

import (
	"context"
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

func logDuration(ctx context.Context, name, method string, startTime time.Duration, err error) {
	duration := timex.Since(startTime)
	logger := logx.WithContext(ctx).WithDuration(duration)
	if err != nil {
		logger.Infof("mongo(%s) - %s - fail(%s)", name, method, err.Error())
	} else {
		logger.Infof("mongo(%s) - %s - ok", name, method)
	}
}
