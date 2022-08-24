package httpc

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/breaker"
)

type (
	// Option is used to customize the *http.Client.
	Option func(r *http.Request) *http.Request
	// RetryFunc return true will retry the request.
	RetryFunc func(resp *http.Response, err error) bool

	// Service represents a remote HTTP service.
	Service interface {
		// Do sends an HTTP request with the given arguments and returns an HTTP response.
		Do(ctx context.Context, method, url string, data interface{}) (*http.Response, error)
		// DoRequest sends a HTTP request to the service.
		DoRequest(r *http.Request) (*http.Response, error)
		// DoRequestWithRetry sends an HTTP request to the service with retry fucntion.
		DoRequestWithRetry(r *http.Request, retryFunc RetryFunc, retryTimes int) (*http.Response, error)
	}

	namedService struct {
		name string
		cli  *http.Client
		opts []Option
	}
)

// NewService returns a remote service with the given name.
// opts are used to customize the *http.Client.
func NewService(name string, opts ...Option) Service {
	return NewServiceWithClient(name, http.DefaultClient, opts...)
}

// NewServiceWithClient returns a remote service with the given name.
// opts are used to customize the *http.Client.
func NewServiceWithClient(name string, cli *http.Client, opts ...Option) Service {
	return namedService{
		name: name,
		cli:  cli,
		opts: opts,
	}
}

// Do sends an HTTP request with the given arguments and returns an HTTP response.
func (s namedService) Do(ctx context.Context, method, url string, data interface{}) (*http.Response, error) {
	req, err := buildRequest(ctx, method, url, data)
	if err != nil {
		return nil, err
	}

	return s.DoRequest(req)
}

// DoRequest sends an HTTP request to the service.
func (s namedService) DoRequest(r *http.Request) (*http.Response, error) {
	return request(r, s)
}

// DoRequestWithRetry sends an HTTP request to the service with retry function.
func (s namedService) DoRequestWithRetry(r *http.Request, retryFunc RetryFunc, retryTimes int) (resp *http.Response, err error) {
	if retryTimes <= 0 {
		retryTimes = 3
	}
	for i := 0; i < retryTimes+1; i++ {
		resp, err = request(cloneReq(r), s)
		if !retryFunc(cloneResp(resp), err) {
			break
		}
	}
	return resp, err
}

func (s namedService) do(r *http.Request) (resp *http.Response, err error) {
	for _, opt := range s.opts {
		r = opt(r)
	}

	brk := breaker.GetBreaker(s.name)
	err = brk.DoWithAcceptable(func() error {
		resp, err = s.cli.Do(r)
		return err
	}, func(err error) bool {
		return err == nil && resp.StatusCode < http.StatusInternalServerError
	})

	return
}

func cloneReq(req *http.Request) *http.Request {
	if req == nil {
		return nil
	} else if req.Body == nil {
		return req
	}
	r := *req
	var b bytes.Buffer
	b.ReadFrom(req.Body)
	req.Body = io.NopCloser(&b)
	r.Body = io.NopCloser(bytes.NewReader(b.Bytes()))
	return &r
}

func cloneResp(resp *http.Response) *http.Response {
	if resp == nil {
		return nil
	} else if resp.Body == nil {
		return resp
	}
	r := *resp
	var b bytes.Buffer
	b.ReadFrom(resp.Body)
	resp.Body = io.NopCloser(&b)
	r.Body = io.NopCloser(bytes.NewReader(b.Bytes()))
	return &r
}
