package serverinterceptors

import (
	"context"
	"encoding/json"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stat"
	"github.com/tal-tech/go-zero/core/timex"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

const serverSlowThreshold = time.Millisecond * 500

func UnaryStatInterceptor(metrics *stat.Metrics) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer handleCrash(func(r interface{}) {
			err = toPanicError(r)
		})

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
	} else if duration > serverSlowThreshold {
		logx.WithContext(ctx).WithDuration(duration).Slowf("[RPC] slowcall - %s - %s - %s",
			addr, method, string(content))
	} else {
		logx.WithContext(ctx).WithDuration(duration).Infof("%s - %s - %s", addr, method, string(content))
	}
}
