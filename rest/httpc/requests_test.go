package httpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()
	_, err := Get("foo", "tcp://bad request")
	assert.NotNil(t, err)
	resp, err := Get("foo", svr.URL)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDoNotFound(t *testing.T) {
	svr := httptest.NewServer(http.NotFoundHandler())
	defer svr.Close()
	_, err := Post("foo", "tcp://bad request", "application/json", nil)
	assert.NotNil(t, err)
	resp, err := Post("foo", svr.URL, "application/json", nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestDoMoved(t *testing.T) {
	svr := httptest.NewServer(http.RedirectHandler("/foo", http.StatusMovedPermanently))
	defer svr.Close()
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)
	_, err = Do("foo", req)
	// too many redirects
	assert.NotNil(t, err)
}
