package httpc

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/rest/internal/header"
)

func TestNamedService_DoRequest(t *testing.T) {
	svr := httptest.NewServer(http.RedirectHandler("/foo", http.StatusMovedPermanently))
	defer svr.Close()
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)
	service := NewService("foo")
	_, err = service.DoRequest(req)
	// too many redirects
	assert.NotNil(t, err)
}

func TestNamedService_DoRequestGet(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("foo", r.Header.Get("foo"))
	}))
	defer svr.Close()
	service := NewService("foo", func(r *http.Request) *http.Request {
		r.Header.Set("foo", "bar")
		return r
	})
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)
	resp, err := service.DoRequest(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "bar", resp.Header.Get("foo"))
}

func TestNamedService_DoRequestPost(t *testing.T) {
	svr := httptest.NewServer(http.NotFoundHandler())
	defer svr.Close()
	service := NewService("foo")
	req, err := http.NewRequest(http.MethodPost, svr.URL, nil)
	assert.Nil(t, err)
	req.Header.Set(header.ContentType, header.JsonContentType)
	resp, err := service.DoRequest(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestNamedService_Do(t *testing.T) {
	type Data struct {
		Key    string `path:"key"`
		Value  int    `form:"value"`
		Header string `header:"X-Header"`
		Body   string `json:"body"`
	}

	svr := httptest.NewServer(http.NotFoundHandler())
	defer svr.Close()

	service := NewService("foo")
	data := Data{
		Key:    "foo",
		Value:  10,
		Header: "my-header",
		Body:   "my body",
	}
	resp, err := service.Do(context.Background(), http.MethodPost, svr.URL+"/nodes/:key", data)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestNamedService_DoBadRequest(t *testing.T) {
	val := struct {
		Value string `json:"value,options=[a,b]"`
	}{
		Value: "c",
	}

	service := NewService("foo")
	_, err := service.Do(context.Background(), http.MethodPost, "/nodes/:key", val)
	assert.NotNil(t, err)
}

func TestNamedService_DoRequestWithRetry(t *testing.T) {
	n := 0
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if n == 2 {
			w.Write([]byte("ok"))
			return
		}
		n++
	}))
	defer svr.Close()
	service := NewService("retry")
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)
	resp, err := service.DoRequestWithRetry(req, func(resp *http.Response, err error) bool {
		if resp != nil && resp.Body != nil {
			bs, _ := io.ReadAll(resp.Body)
			return string(bs) != "ok"
		}
		return true
	}, 5)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, n)
	bs, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "ok", string(bs))
}

func TestNamedService_DoRequestPostWithRetry(t *testing.T) {
	n := 0
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if n == 2 {
			w.Write([]byte("ok"))
			return
		}
		n++
	}))
	defer svr.Close()
	service := NewService("retry")
	req, err := http.NewRequest(http.MethodPost, svr.URL, strings.NewReader(`{}`))
	assert.Nil(t, err)
	resp, err := service.DoRequestWithRetry(req, func(resp *http.Response, err error) bool {
		if resp != nil && resp.Body != nil {
			bs, _ := io.ReadAll(resp.Body)
			return string(bs) != "ok"
		}
		return true
	}, 5)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, n)
	bs, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "ok", string(bs))
}
