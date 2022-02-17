package clientinterceptors

import (
	"context"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/prometheus"
	"github.com/zeromicro/go-zero/core/timex"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const clientNamespace = "rpc_client"

var (
	metricClientReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: clientNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "rpc client requests duration(ms).",
		Labels:    []string{"method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	metricClientReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: clientNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "rpc client requests code count.",
		Labels:    []string{"method", "code"},
	})
)

// PrometheusInterceptor is an interceptor that reports to prometheus server.
func PrometheusInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if !prometheus.Enabled() {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	startTime := timex.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	metricClientReqDur.Observe(int64(timex.Since(startTime)/time.Millisecond), method)
	metricClientReqCodeTotal.Inc(method, strconv.Itoa(int(status.Code(err))))
	return err
}
