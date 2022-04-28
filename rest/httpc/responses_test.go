package httpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/rest/internal/header"
)

func TestParse(t *testing.T) {
	var val struct {
		Foo   string `header:"foo"`
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("foo", "bar")
		w.Header().Set(header.ContentType, header.JsonContentType)
		w.Write([]byte(`{"name":"kevin","value":100}`))
	}))
	defer svr.Close()
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)
	resp, err := DoRequest(req)
	assert.Nil(t, err)
	assert.Nil(t, Parse(resp, &val))
	assert.Equal(t, "bar", val.Foo)
	assert.Equal(t, "kevin", val.Name)
	assert.Equal(t, 100, val.Value)
}

func TestParseHeaderError(t *testing.T) {
	var val struct {
		Foo int `header:"foo"`
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("foo", "bar")
		w.Header().Set(header.ContentType, header.JsonContentType)
	}))
	defer svr.Close()
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)
	resp, err := DoRequest(req)
	assert.Nil(t, err)
	assert.NotNil(t, Parse(resp, &val))
}

func TestParseNoBody(t *testing.T) {
	var val struct {
		Foo string `header:"foo"`
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("foo", "bar")
		w.Header().Set(header.ContentType, header.JsonContentType)
	}))
	defer svr.Close()
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)
	resp, err := DoRequest(req)
	assert.Nil(t, err)
	assert.Nil(t, Parse(resp, &val))
	assert.Equal(t, "bar", val.Foo)
}
