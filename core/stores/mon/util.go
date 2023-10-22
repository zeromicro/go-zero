package mon

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
)

const mongoAddrSep = ","

var errPlaceholder = errors.New("placeholder")

// FormatAddr formats mongo hosts to a string.
func FormatAddr(hosts []string) string {
	return strings.Join(hosts, mongoAddrSep)
}

func logDuration(ctx context.Context, name, method string, startTime time.Duration, err error) {
	logDurationWithDoc(ctx, name, method, startTime, err)
}

func logDurationWithDoc(ctx context.Context, name, method string,
	startTime time.Duration, err error, docs ...any) {
	duration := timex.Since(startTime)
	logger := logx.WithContext(ctx).WithDuration(duration)
	var content []byte
	jerr := errPlaceholder
	if len(docs) > 0 {
		content, jerr = json.Marshal(docs)
	}

	if err == nil {
		// jerr should not be non-nil, but we don't care much on this,
		// if non-nil, we just log without docs.
		if jerr != nil {
			if logSlowMon.True() && duration > slowThreshold.Load() {
				logger.Slowf("mongo(%s) - slowcall - %s - ok", name, method)
			} else if logMon.True() {
				logger.Infof("mongo(%s) - %s - ok", name, method)
			}
		} else {
			if logSlowMon.True() && duration > slowThreshold.Load() {
				logger.Slowf("mongo(%s) - slowcall - %s - ok - %s",
					name, method, string(content))
			} else if logMon.True() {
				logger.Infof("mongo(%s) - %s - ok - %s",
					name, method, string(content))
			}
		}

		return
	}

	if jerr != nil {
		logger.Errorf("mongo(%s) - %s - fail(%s)", name, method, err.Error())
	} else {
		logger.Errorf("mongo(%s) - %s - fail(%s) - %s",
			name, method, err.Error(), string(content))
	}
}
