package httpc

import (
	"io"
	"net/http"
)

// Do sends an HTTP request to the service assocated with the given key.
func Do(key string, r *http.Request, opts ...Option) (*http.Response, error) {
	return NewService(key, opts...).Do(r)
}

// Get sends an HTTP GET request to the service assocated with the given key.
func Get(key, url string, opts ...Option) (*http.Response, error) {
	return NewService(key, opts...).Get(url)
}

// Post sends an HTTP POST request to the service assocated with the given key.
func Post(key, url, contentType string, body io.Reader, opts ...Option) (*http.Response, error) {
	return NewService(key, opts...).Post(url, contentType, body)
}
