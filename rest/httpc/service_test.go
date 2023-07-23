package httpc

import (
	"context"
	"net/http"
	"net/http/httptest"
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
