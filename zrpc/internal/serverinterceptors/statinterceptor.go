package serverinterceptors

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/timex"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

const defaultSlowThreshold = time.Millisecond * 500

var slowThreshold = syncx.ForAtomicDuration(defaultSlowThreshold)

// SetSlowThreshold sets the slow threshold.
func SetSlowThreshold(threshold time.Duration) {
	slowThreshold.Set(threshold)
}

// UnaryStatInterceptor returns a func that uses given metrics to report stats.
func UnaryStatInterceptor(metrics *stat.Metrics) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			metrics.Add(stat.Task{
				Duration: duration,
			})
			logDuration(ctx, info.FullMethod, req, duration)
		}()

		return handler(ctx, req)
	}
}

func logDuration(ctx context.Context, method string, req interface{}, duration time.Duration) {
	var addr string
	client, ok := peer.FromContext(ctx)
	if ok {
		addr = client.Addr.String()
	}
	content, err := json.Marshal(req)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s - %s", addr, err.Error())
	} else if duration > slowThreshold.Load() {
		logx.WithContext(ctx).WithDuration(duration).Slowf("[RPC] slowcall - %s - %s - %s",
			addr, method, string(content))
	} else {
		logx.WithContext(ctx).WithDuration(duration).Infof("%s - %s - %s", addr, method, string(content))
	}
}
