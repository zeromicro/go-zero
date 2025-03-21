package clientinterceptors

import (
	"context"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/timex"
)

const clientNamespace = "rpc_client"

var (
	defaultDurationBuckets   = []float64{1, 2, 5, 10, 25, 50, 100, 250, 500, 1000, 2000, 5000}
	metricClientReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: clientNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "rpc client requests code count.",
		Labels:    []string{"method", "code"},
	})
)

// PrometheusInterceptor is an interceptor that reports to prometheus server.
func PrometheusInterceptor(buckets []float64) grpc.UnaryClientInterceptor {
	if len(buckets) == 0 {
		buckets = defaultDurationBuckets
	}

	metricDurationHistogram := initMetricClientReqDur(buckets)

	return func(ctx context.Context, method string, req, reply any,
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		startTime := timex.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		metricDurationHistogram.Observe(timex.Since(startTime).Milliseconds(), method)
		metricClientReqCodeTotal.Inc(method, strconv.Itoa(int(status.Code(err))))
		return err
	}
}

func initMetricClientReqDur(buckets []float64) metric.HistogramVec {
	return metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: clientNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "rpc client requests duration(ms).",
		Labels:    []string{"method"},
		Buckets:   buckets,
	})
}
