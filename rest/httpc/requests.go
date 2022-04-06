package httpc

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpc/internal"
)

var interceptors = []internal.Interceptor{
	internal.LogInterceptor,
}

// Do sends an HTTP request with the given arguments and returns an HTTP response.
// data is automatically marshal into a *httpRequest, typically it's defined in an API file.
func Do(ctx context.Context, method, url string, data interface{}) (*http.Response, error) {
	req, err := buildRequest(ctx, method, url, data)
	if err != nil {
		return nil, err
	}

	return DoRequest(req)
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

func buildRequest(ctx context.Context, method, url string, data interface{}) (*http.Request, error) {
	val, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return http.NewRequestWithContext(ctx, method, url, bytes.NewReader(val))
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
