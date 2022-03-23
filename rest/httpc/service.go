package httpc

import (
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc/internal"
)

var interceptors = []internal.Interceptor{
	internal.LogInterceptor,
}

type (
	// Option is used to customize the *http.Client.
	Option func(r *http.Request) *http.Request

	// Service represents a remote HTTP service.
	Service interface {
		// Do sends an HTTP request to the service.
		Do(r *http.Request) (*http.Response, error)
		// Get sends an HTTP GET request to the service.
		Get(url string) (*http.Response, error)
		// Post sends an HTTP POST request to the service.
		Post(url, contentType string, body io.Reader) (*http.Response, error)
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

// Do sends an HTTP request to the service.
func (s namedService) Do(r *http.Request) (resp *http.Response, err error) {
	var respHandlers []internal.ResponseHandler
	for _, interceptor := range interceptors {
		var h internal.ResponseHandler
		r, h = interceptor(r)
		respHandlers = append(respHandlers, h)
	}

	resp, err = s.doRequest(r)
	if err != nil {
		logx.Errorf("[HTTP] %s %s/%s - %v", r.Method, r.Host, r.RequestURI, err)
		return
	}

	for i := len(respHandlers) - 1; i >= 0; i-- {
		respHandlers[i](resp)
	}

	return
}

// Get sends an HTTP GET request to the service.
func (s namedService) Get(url string) (*http.Response, error) {
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return s.Do(r)
}

// Post sends an HTTP POST request to the service.
func (s namedService) Post(url, ctype string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	r.Header.Set(contentType, ctype)
	return s.Do(r)
}

func (s namedService) doRequest(r *http.Request) (resp *http.Response, err error) {
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
