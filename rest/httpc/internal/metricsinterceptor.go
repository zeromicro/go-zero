package internal

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/timex"
)

const clientNamespace = "httpc_client"

var (
	MetricClientReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: clientNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http client requests duration(ms).",
		Labels:    []string{"name", "method", "url"},
		Buckets:   []float64{0.25, 0.5, 1, 2, 5, 10, 25, 50, 100, 250, 500, 1000, 2000, 5000, 10000, 15000},
	})

	MetricClientReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: clientNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http client requests code count.",
		Labels:    []string{"name", "method", "url", "code"},
	})
)

type MetricsURLRewriter func(u url.URL) string

func MetricsInterceptor(name string, pr MetricsURLRewriter) Interceptor {
	return func(r *http.Request) (*http.Request, ResponseHandler) {
		startTime := timex.Now()
		return r, func(resp *http.Response, err error) {
			var code int
			var path string

			// error or resp is nil, set code=500
			if err != nil || resp == nil {
				code = http.StatusInternalServerError
			} else {
				code = resp.StatusCode
			}

			u := cleanURL(*r.URL)
			method := r.Method
			if pr != nil {
				path = pr(u)
			} else {
				path = u.String()
			}

			MetricClientReqDur.ObserveFloat(float64(timex.Since(startTime))/float64(time.Millisecond), name, method, path)
			MetricClientReqCodeTotal.Inc(name, method, path, strconv.Itoa(code))
		}
	}
}

func cleanURL(r url.URL) url.URL {
	r.RawQuery = ""
	r.RawFragment = ""
	r.User = nil
	return r
}
