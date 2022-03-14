package httpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNamedService_Do(t *testing.T) {
	svr := httptest.NewServer(http.RedirectHandler("/foo", http.StatusMovedPermanently))
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)
	service := NewService("foo")
	_, err = service.Do(req)
	// too many redirects
	assert.NotNil(t, err)
}

func TestNamedService_Get(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	service := NewService("foo")
	resp, err := service.Get(svr.URL, func(cli *http.Client) {
		cli.Transport = http.DefaultTransport
	})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestNamedService_Post(t *testing.T) {
	svr := httptest.NewServer(http.NotFoundHandler())
	service := NewService("foo")
	_, err := service.Post("tcp://bad request", "application/json", nil)
	assert.NotNil(t, err)
	resp, err := service.Post(svr.URL, "application/json", nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
