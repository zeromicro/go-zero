package httpc

import (
	"io"
	"net/http"
)

type (
	Service interface {
		Do(r *http.Request, opts ...Option) (*http.Response, error)
		Get(url string, opts ...Option) (*http.Response, error)
		Post(url, contentType string, body io.Reader, opts ...Option) (*http.Response, error)
	}

	namedService struct {
		name string
		opts []Option
	}
)

func NewService(name string, opts ...Option) Service {
	return namedService{
		name: name,
		opts: opts,
	}
}

func (s namedService) Do(r *http.Request, opts ...Option) (*http.Response, error) {
	return Do(s.name, r, append(s.opts, opts...)...)
}

func (s namedService) Get(url string, opts ...Option) (*http.Response, error) {
	return Get(s.name, url, append(s.opts, opts...)...)
}

func (s namedService) Post(url, contentType string, body io.Reader, opts ...Option) (
	*http.Response, error) {
	return Post(s.name, url, contentType, body, append(s.opts, opts...)...)
}
