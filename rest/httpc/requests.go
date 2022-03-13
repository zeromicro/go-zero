package httpc

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
	"go.opentelemetry.io/otel/propagation"
)

func Do(key string, r *http.Request) (resp *http.Response, err error) {
	brk := breaker.GetBreaker(key)
	start := timex.Now()
	err = brk.DoWithAcceptable(func() error {
		resp, err = http.DefaultClient.Do(r)
		return err
	}, func(err error) bool {
		return err == nil && resp.StatusCode < http.StatusInternalServerError
	})

	if err != nil {
		logx.Errorf("[HTTP] %s %s/%s - %v", r.Method, r.Host, r.RequestURI, err)
		return
	}

	duration := timex.Since(start)
	var tc propagation.TraceContext
	ctx := tc.Extract(context.Background(), propagation.HeaderCarrier(resp.Header))
	logger := logx.WithContext(ctx).WithDuration(duration)
	if isOkResponse(resp.StatusCode) {
		logger.Infof("[HTTP] %d - %s %s/%s", resp.StatusCode, r.Method, r.Host, r.RequestURI)
	} else {
		logger.Errorf("[HTTP] %d - %s %s/%s", resp.StatusCode, r.Method, r.Host, r.RequestURI)
	}

	return
}

func isOkResponse(code int) bool {
	return code < http.StatusBadRequest
}
