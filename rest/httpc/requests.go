package httpc

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc/internal"
)

var interceptors = []internal.Interceptor{
	internal.LogInterceptor,
}

type Option func(cli *http.Client)

// Do sends an HTTP request to the service assocated with the given key.
func Do(key string, r *http.Request, opts ...Option) (resp *http.Response, err error) {
	var respHandlers []internal.ResponseHandler
	for _, interceptor := range interceptors {
		var h internal.ResponseHandler
		r, h = interceptor(r)
		respHandlers = append(respHandlers, h)
	}

	resp, err = doRequest(key, r, opts...)
	if err != nil {
		logx.Errorf("[HTTP] %s %s/%s - %v", r.Method, r.Host, r.RequestURI, err)
		return
	}

	for i := len(respHandlers) - 1; i >= 0; i-- {
		respHandlers[i](resp)
	}

	return
}

func doRequest(key string, r *http.Request, opts ...Option) (resp *http.Response, err error) {
	brk := breaker.GetBreaker(key)
	err = brk.DoWithAcceptable(func() error {
		var cli http.Client
		for _, opt := range opts {
			opt(&cli)
		}
		resp, err = cli.Do(r)
		return err
	}, func(err error) bool {
		return err == nil && resp.StatusCode < http.StatusInternalServerError
	})

	return
}
