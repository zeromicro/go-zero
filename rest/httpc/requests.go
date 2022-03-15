package httpc

import (
	"io"
	"net/http"
)

// Do sends an HTTP request to the service assocated with the given key.
func Do(key string, r *http.Request) (*http.Response, error) {
	return NewService(key).Do(r)
}

// Get sends an HTTP GET request to the service assocated with the given key.
func Get(key, url string) (*http.Response, error) {
	return NewService(key).Get(url)
}

// Post sends an HTTP POST request to the service assocated with the given key.
func Post(key, url, contentType string, body io.Reader) (*http.Response, error) {
	return NewService(key).Post(url, contentType, body)
}
