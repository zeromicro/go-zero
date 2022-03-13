package internal

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
	"go.opentelemetry.io/otel/propagation"
)

func LogInterceptor(r *http.Request) (*http.Request, ResponseHandler) {
	start := timex.Now()
	return r, func(resp *http.Response) {
		duration := timex.Since(start)
		var tc propagation.TraceContext
		ctx := tc.Extract(r.Context(), propagation.HeaderCarrier(resp.Header))
		logger := logx.WithContext(ctx).WithDuration(duration)
		if isOkResponse(resp.StatusCode) {
			logger.Infof("[HTTP] %d - %s %s/%s", resp.StatusCode, r.Method, r.Host, r.RequestURI)
		} else {
			logger.Errorf("[HTTP] %d - %s %s/%s", resp.StatusCode, r.Method, r.Host, r.RequestURI)
		}
	}
}

func isOkResponse(code int) bool {
	return code < http.StatusBadRequest
}
