package handler

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/timex"
	"github.com/zeromicro/go-zero/rest/internal/response"
)

const serverNamespace = "http_server"

var (
	serverReqDurBuckets    = []float64{5, 10, 25, 50, 100, 250, 500, 750, 1000}
	metricServerReqDurOnce sync.Once

	metricServerReqDur metric.HistogramVec

	metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http server requests error count.",
		Labels:    []string{"path", "method", "code"},
	})
)

// SetServerReqDurBuckets sets buckets for rest server requests duration.
// It must be called before PrometheusHandler is used.
func SetServerReqDurBuckets(buckets []float64) {
	serverReqDurBuckets = buckets
}

// PrometheusHandler returns a middleware that reports stats to prometheus.
func PrometheusHandler(path, method string) func(http.Handler) http.Handler {
	initMetricServerReqDur()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := timex.Now()
			cw := response.NewWithCodeResponseWriter(w)
			defer func() {
				code := strconv.Itoa(cw.Code)
				metricServerReqDur.Observe(timex.Since(startTime).Milliseconds(), path, method, code)
				metricServerReqCodeTotal.Inc(path, method, code)
			}()

			next.ServeHTTP(cw, r)
		})
	}
}

func initMetricServerReqDur() {
	metricServerReqDurOnce.Do(func() {
		metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
			Namespace: serverNamespace,
			Subsystem: "requests",
			Name:      "duration_ms",
			Help:      "http server requests duration(ms).",
			Labels:    []string{"path", "method", "code"},
			Buckets:   serverReqDurBuckets,
		})
	})
}
