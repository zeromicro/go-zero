package clientinterceptors

import (
	"context"
	"path"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/timex"
	"google.golang.org/grpc"
)

const defaultSlowThreshold = time.Millisecond * 500

var (
	notLoggingContentMethods sync.Map
	slowThreshold            = syncx.ForAtomicDuration(defaultSlowThreshold)
)

// DurationInterceptor is an interceptor that logs the processing time.
func DurationInterceptor(ctx context.Context, method string, req, reply any,
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	serverName := path.Join(cc.Target(), method)
	start := timex.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		logger := logx.WithContext(ctx).WithDuration(timex.Since(start))
		_, ok := notLoggingContentMethods.Load(method)
		if ok {
			logger.Errorf("fail - %s - %s", serverName, err.Error())
		} else {
			logger.Errorf("fail - %s - %v - %s", serverName, req, err.Error())
		}
	} else {
		elapsed := timex.Since(start)
		if elapsed > slowThreshold.Load() {
			logger := logx.WithContext(ctx).WithDuration(elapsed)
			_, ok := notLoggingContentMethods.Load(method)
			if ok {
				logger.Slowf("[RPC] ok - slowcall - %s", serverName)
			} else {
				logger.Slowf("[RPC] ok - slowcall - %s - %v - %v", serverName, req, reply)
			}
		}
	}

	return err
}

// DontLogContentForMethod disable logging content for given method.
func DontLogContentForMethod(method string) {
	notLoggingContentMethods.Store(method, lang.Placeholder)
}

// SetSlowThreshold sets the slow threshold.
func SetSlowThreshold(threshold time.Duration) {
	slowThreshold.Set(threshold)
}
