package httpc

import (
	"errors"
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
		w.Header().Set(header.ContentType, header.ContentTypeJson)
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
		w.Header().Set(header.ContentType, header.ContentTypeJson)
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
		w.Header().Set(header.ContentType, header.ContentTypeJson)
	}))
	defer svr.Close()
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)
	resp, err := DoRequest(req)
	assert.Nil(t, err)
	assert.Nil(t, Parse(resp, &val))
	assert.Equal(t, "bar", val.Foo)
}

func TestParseWithZeroValue(t *testing.T) {
	var val struct {
		Foo int `header:"foo"`
		Bar int `json:"bar"`
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("foo", "0")
		w.Header().Set(header.ContentType, header.ContentTypeJson)
		w.Write([]byte(`{"bar":0}`))
	}))
	defer svr.Close()
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)
	resp, err := DoRequest(req)
	assert.Nil(t, err)
	assert.Nil(t, Parse(resp, &val))
	assert.Equal(t, 0, val.Foo)
	assert.Equal(t, 0, val.Bar)
}

func TestParseWithNegativeContentLength(t *testing.T) {
	var val struct {
		Bar int `json:"bar"`
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(header.ContentType, header.ContentTypeJson)
		w.Write([]byte(`{"bar":0}`))
	}))
	defer svr.Close()
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)

	tests := []struct {
		name   string
		length int64
	}{
		{
			name:   "negative",
			length: -1,
		},
		{
			name:   "zero",
			length: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := DoRequest(req)
			resp.ContentLength = test.length
			assert.Nil(t, err)
			assert.Nil(t, Parse(resp, &val))
			assert.Equal(t, 0, val.Bar)
		})
	}
}

func TestParseWithNegativeContentLengthNoBody(t *testing.T) {
	var val struct{}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(header.ContentType, header.ContentTypeJson)
	}))
	defer svr.Close()
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)

	tests := []struct {
		name   string
		length int64
	}{
		{
			name:   "negative",
			length: -1,
		},
		{
			name:   "zero",
			length: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := DoRequest(req)
			resp.ContentLength = test.length
			assert.Nil(t, err)
			assert.Nil(t, Parse(resp, &val))
		})
	}
}

func TestParseJsonBody_BodyError(t *testing.T) {
	var val struct{}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(header.ContentType, header.ContentTypeJson)
	}))
	defer svr.Close()
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)

	resp, err := DoRequest(req)
	resp.ContentLength = -1
	resp.Body = mockedReader{}
	assert.Nil(t, err)
	assert.NotNil(t, Parse(resp, &val))
}

type mockedReader struct{}

func (m mockedReader) Close() error {
	return nil
}

func (m mockedReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("dummy")
}
