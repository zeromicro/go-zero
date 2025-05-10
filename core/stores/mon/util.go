package mon

import (
	"context"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/jsonx"
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
		logger.Errorf("mongo(%s) - %s - fail(%s)", name, method, err.Error())
		return
	}

	if logSlowMon.True() && duration > slowThreshold.Load() {
		logger.Slowf("[MONGO] mongo(%s) - slowcall - %s - ok", name, method)
	} else if logMon.True() {
		logger.Infof("mongo(%s) - %s - ok", name, method)
	}
}

func logDurationWithDocs(ctx context.Context, name, method string, startTime time.Duration,
	err error, docs ...any) {
	duration := timex.Since(startTime)
	logger := logx.WithContext(ctx).WithDuration(duration)

	content, jerr := jsonx.MarshalToString(docs)
	// jerr should not be non-nil, but we don't care much on this,
	// if non-nil, we just log without docs.
	if jerr != nil {
		if err != nil {
			logger.Errorf("mongo(%s) - %s - fail(%s)", name, method, err.Error())
		} else if logSlowMon.True() && duration > slowThreshold.Load() {
			logger.Slowf("[MONGO] mongo(%s) - slowcall - %s - ok", name, method)
		} else if logMon.True() {
			logger.Infof("mongo(%s) - %s - ok", name, method)
		}
		return
	}

	if err != nil {
		logger.Errorf("mongo(%s) - %s - fail(%s) - %s", name, method, err.Error(), content)
	} else if logSlowMon.True() && duration > slowThreshold.Load() {
		logger.Slowf("[MONGO] mongo(%s) - slowcall - %s - ok - %s", name, method, content)
	} else if logMon.True() {
		logger.Infof("mongo(%s) - %s - ok - %s", name, method, content)
	}
}
