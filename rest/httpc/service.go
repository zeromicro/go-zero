package httpc

import (
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc/internal"
)

// ContentType means Content-Type.
const ContentType = "Content-Type"

var interceptors = []internal.Interceptor{
	internal.LogInterceptor,
}

type (
	// Option is used to customize the *http.Client.
	Option func(cli *http.Client)

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
	}
)

// NewService returns a remote service with the given name.
// opts are used to customize the *http.Client.
func NewService(name string, opts ...Option) Service {
	var cli *http.Client

	if len(opts) == 0 {
		cli = http.DefaultClient
	} else {
		cli = &http.Client{}
		for _, opt := range opts {
			opt(cli)
		}
	}

	return namedService{
		name: name,
		cli:  cli,
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
func (s namedService) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	r.Header.Set(ContentType, contentType)
	return s.Do(r)
}

func (s namedService) doRequest(r *http.Request) (resp *http.Response, err error) {
	brk := breaker.GetBreaker(s.name)
	err = brk.DoWithAcceptable(func() error {
		resp, err = s.cli.Do(r)
		return err
	}, func(err error) bool {
		return err == nil && resp.StatusCode < http.StatusInternalServerError
	})

	return
}
