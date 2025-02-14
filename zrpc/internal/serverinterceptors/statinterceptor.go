package serverinterceptors

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/timex"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

const defaultSlowThreshold = time.Millisecond * 500

var (
	ignoreContentMethods sync.Map
	slowThreshold        = syncx.ForAtomicDuration(defaultSlowThreshold)
)

// StatConf defines the static configuration for stat interceptor.
type StatConf struct {
	SlowThreshold        time.Duration `json:",default=500ms"`
	IgnoreContentMethods []string      `json:",optional"`
}

// DontLogContentForMethod disable logging content for given method.
// Deprecated: use StatConf instead.
func DontLogContentForMethod(method string) {
	ignoreContentMethods.Store(method, lang.Placeholder)
}

// SetSlowThreshold sets the slow threshold.
// Deprecated: use StatConf instead.
func SetSlowThreshold(threshold time.Duration) {
	slowThreshold.Set(threshold)
}

// UnaryStatInterceptor returns a func that uses given metrics to report stats.
func UnaryStatInterceptor(metrics *stat.Metrics, conf StatConf) grpc.UnaryServerInterceptor {
	staticNotLoggingContentMethods := collection.NewSet()
	staticNotLoggingContentMethods.AddStr(conf.IgnoreContentMethods...)

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			metrics.Add(stat.Task{
				Duration: duration,
			})
			logDuration(ctx, info.FullMethod, req, duration,
				staticNotLoggingContentMethods, conf.SlowThreshold)
		}()

		return handler(ctx, req)
	}
}

func isSlow(duration, durationThreshold time.Duration) bool {
	return duration > slowThreshold.Load() ||
		(durationThreshold > 0 && duration > durationThreshold)
}

func logDuration(ctx context.Context, method string, req any, duration time.Duration,
	ignoreMethods *collection.Set, durationThreshold time.Duration) {
	var addr string
	client, ok := peer.FromContext(ctx)
	if ok {
		addr = client.Addr.String()
	}

	logger := logx.WithContext(ctx).WithDuration(duration)
	if !shouldLogContent(method, ignoreMethods) {
		if isSlow(duration, durationThreshold) {
			logger.Slowf("[RPC] slowcall - %s - %s", addr, method)
		}
	} else {
		content, err := json.Marshal(req)
		if err != nil {
			logx.WithContext(ctx).Errorf("%s - %s", addr, err.Error())
		} else if isSlow(duration, durationThreshold) {
			logger.Slowf("[RPC] slowcall - %s - %s - %s", addr, method, string(content))
		} else {
			logger.Infof("%s - %s - %s", addr, method, string(content))
		}
	}
}

func shouldLogContent(method string, ignoreMethods *collection.Set) bool {
	_, ok := ignoreContentMethods.Load(method)
	return !ok && !ignoreMethods.Contains(method)
}
