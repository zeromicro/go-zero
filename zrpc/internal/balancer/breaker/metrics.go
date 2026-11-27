package breaker

import "github.com/zeromicro/go-zero/core/metric"

var metricInstanceBreakerTriggered = metric.NewCounterVec(&metric.CounterVecOpts{
	Namespace: "rpc_client",
	Subsystem: "breaker",
	Name:      "instance_triggered_total",
	Help:      "Total number of requests rejected by instance breaker",
	Labels:    []string{"addr", "method"},
})
