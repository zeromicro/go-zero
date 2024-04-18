package httpc

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/breaker"
)

type (
	// Option is used to customize the *http.Client.
	Option func(r *http.Request) *http.Request

	// Service represents a remote HTTP service.
	Service interface {
		// Do sends an HTTP request with the given arguments and returns an HTTP response.
		Do(ctx context.Context, method, url string, data any) (*http.Response, error)
		// DoRequest sends a HTTP request to the service.
		DoRequest(r *http.Request) (*http.Response, error)
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
func (s namedService) Do(ctx context.Context, method, url string, data any) (*http.Response, error) {
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

func (s namedService) do(r *http.Request) (resp *http.Response, err error) {
	for _, opt := range s.opts {
		r = opt(r)
	}

	brk := breaker.GetBreaker(s.name)
	err = brk.DoWithAcceptableCtx(r.Context(), func() error {
		resp, err = s.cli.Do(r)
		return err
	}, func(err error) bool {
		return err == nil && resp.StatusCode < http.StatusInternalServerError
	})

	return
}
