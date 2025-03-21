package serverinterceptors

import (
	"context"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/timex"
)

const serverNamespace = "rpc_server"

var (
	defaultDurationBuckets   = []float64{1, 2, 5, 10, 25, 50, 100, 250, 500, 1000, 2000, 5000}
	metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "rpc server requests code count.",
		Labels:    []string{"method", "code"},
	})
)

// UnaryPrometheusInterceptor reports the statistics to the prometheus server.
func UnaryPrometheusInterceptor(buckets []float64) grpc.UnaryServerInterceptor {
	if len(buckets) == 0 {
		buckets = defaultDurationBuckets
	}

	metricDurationHistogram := initMetricServerReqDur(buckets)

	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startTime := timex.Now()
		resp, err := handler(ctx, req)
		metricDurationHistogram.Observe(timex.Since(startTime).Milliseconds(), info.FullMethod)
		metricServerReqCodeTotal.Inc(info.FullMethod, strconv.Itoa(int(status.Code(err))))
		return resp, err
	}
}

func initMetricServerReqDur(buckets []float64) metric.HistogramVec {
	return metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "rpc server requests duration(ms).",
		Labels:    []string{"method"},
		Buckets:   buckets,
	})
}
