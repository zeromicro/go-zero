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
	defaultBuckets = []float64{1, 2, 5, 10, 25, 50, 100, 250, 500, 1000, 2000, 5000}

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

// PrometheusHandler returns a middleware that reports stats to prometheus.
func PrometheusHandler(path, method string, buckets []float64) func(http.Handler) http.Handler {
	if len(buckets) == 0 {
		buckets = defaultBuckets
	}

	initMetricServerReqDur(buckets)

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

func initMetricServerReqDur(buckets []float64) {
	metricServerReqDurOnce.Do(func() {
		metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
			Namespace: serverNamespace,
			Subsystem: "requests",
			Name:      "duration_ms",
			Help:      "http server requests duration(ms).",
			Labels:    []string{"path", "method", "code"},
			Buckets:   buckets,
		})
	})
}
