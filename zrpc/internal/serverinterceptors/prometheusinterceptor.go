package serverinterceptors

import (
	"context"
	"strconv"
	"sync"

	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/timex"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const serverNamespace = "rpc_server"

var (
	rpcServerReqDurBuckets = []float64{1, 2, 5, 10, 25, 50, 100, 250, 500, 1000, 2000, 5000}
	metricServerReqDurOnce sync.Once

	metricServerReqDur metric.HistogramVec

	metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "rpc server requests code count.",
		Labels:    []string{"method", "code"},
	})
)

// SetRpcServerReqDurBuckets sets buckets for rpc server requests duration.
// It must be called before UnaryPrometheusInterceptor is used.
func SetRpcServerReqDurBuckets(buckets []float64) {
	rpcServerReqDurBuckets = buckets
}

// UnaryPrometheusInterceptor reports the statistics to the prometheus server.
func UnaryPrometheusInterceptor(ctx context.Context, req any,
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	initMetricServerReqDur()

	startTime := timex.Now()
	resp, err := handler(ctx, req)
	metricServerReqDur.Observe(timex.Since(startTime).Milliseconds(), info.FullMethod)
	metricServerReqCodeTotal.Inc(info.FullMethod, strconv.Itoa(int(status.Code(err))))
	return resp, err
}

func initMetricServerReqDur() {
	metricServerReqDurOnce.Do(func() {
		metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
			Namespace: serverNamespace,
			Subsystem: "requests",
			Name:      "duration_ms",
			Help:      "rpc server requests duration(ms).",
			Labels:    []string{"method"},
			Buckets:   rpcServerReqDurBuckets,
		})
	})
}
