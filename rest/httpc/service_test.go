package httpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNamedService_Do(t *testing.T) {
	svr := httptest.NewServer(http.RedirectHandler("/foo", http.StatusMovedPermanently))
	defer svr.Close()
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)
	service := NewService("foo")
	_, err = service.Do(req)
	// too many redirects
	assert.NotNil(t, err)
}

func TestNamedService_Get(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("foo", r.Header.Get("foo"))
	}))
	defer svr.Close()
	service := NewService("foo", func(r *http.Request) *http.Request {
		r.Header.Set("foo", "bar")
		return r
	})
	resp, err := service.Get(svr.URL)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "bar", resp.Header.Get("foo"))
}

func TestNamedService_Post(t *testing.T) {
	svr := httptest.NewServer(http.NotFoundHandler())
	defer svr.Close()
	service := NewService("foo")
	_, err := service.Post("tcp://bad request", "application/json", nil)
	assert.NotNil(t, err)
	resp, err := service.Post(svr.URL, "application/json", nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
