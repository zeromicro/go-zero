package httpc

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"

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
		return acceptable(resp, err)
	})

	return
}

// acceptable determines whether the HTTP request/response should be considered
// successful for circuit breaker purposes.
//
// Returns true (acceptable) for:
//   - HTTP status codes < 500 (2xx, 3xx, 4xx)
//   - Context cancellation (user-initiated)
//   - Non-network errors (application-level errors)
//
// Returns false (not acceptable, triggers breaker) for:
//   - HTTP status codes >= 500 (server errors)
//   - context.DeadlineExceeded (timeout)
//   - Network errors (connection refused, DNS failures, etc.)
func acceptable(resp *http.Response, err error) bool {
	if err == nil {
		return resp.StatusCode < http.StatusInternalServerError
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	if errors.Is(err, context.Canceled) {
		return true
	}

	// Unwrap url.Error if present
	var ue *url.Error
	if errors.As(err, &ue) {
		err = ue.Unwrap()
	}

	// Network errors are not acceptable
	var ne net.Error
	return !errors.As(err, &ne)
}
