package httpc

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpc/internal"
)

var interceptors = []internal.Interceptor{
	internal.LogInterceptor,
}

// DoRequest sends an HTTP request and returns an HTTP response.
func DoRequest(r *http.Request) (*http.Response, error) {
	return request(r, defaultClient{})
}

type (
	client interface {
		do(r *http.Request) (*http.Response, error)
	}

	defaultClient struct{}
)

func (c defaultClient) do(r *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(r)
}

func request(r *http.Request, cli client) (*http.Response, error) {
	var respHandlers []internal.ResponseHandler
	for _, interceptor := range interceptors {
		var h internal.ResponseHandler
		r, h = interceptor(r)
		respHandlers = append(respHandlers, h)
	}

	resp, err := cli.do(r)
	for i := len(respHandlers) - 1; i >= 0; i-- {
		respHandlers[i](resp, err)
	}

	return resp, err
}
