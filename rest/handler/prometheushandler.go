package handler

import (
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/timex"
	"github.com/zeromicro/go-zero/rest/internal/response"
)

const serverNamespace = "http_server"

var (
	metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http server requests duration(ms).",
		Labels:    []string{"path", "method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http server requests error count.",
		Labels:    []string{"path", "code", "method"},
	})
)

// PrometheusHandler returns a middleware that reports stats to prometheus.
func PrometheusHandler(path, method string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := timex.Now()
			cw := response.NewWithCodeResponseWriter(w)
			defer func() {
				metricServerReqDur.Observe(timex.Since(startTime).Milliseconds(), path, method)
				metricServerReqCodeTotal.Inc(path, strconv.Itoa(cw.Code), method)
			}()

			next.ServeHTTP(cw, r)
		})
	}
}
