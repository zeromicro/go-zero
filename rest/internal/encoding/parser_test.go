package encoding

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHeaders(t *testing.T) {
	var val struct {
		Foo string `header:"foo"`
		Baz int    `header:"baz"`
		Qux bool   `header:"qux,default=true"`
	}
	r := httptest.NewRequest(http.MethodGet, "/any", http.NoBody)
	r.Header.Set("foo", "bar")
	r.Header.Set("baz", "1")
	assert.Nil(t, ParseHeaders(r.Header, &val))
	assert.Equal(t, "bar", val.Foo)
	assert.Equal(t, 1, val.Baz)
	assert.True(t, val.Qux)
}

func TestParseHeadersMulti(t *testing.T) {
	var val struct {
		Foo []string `header:"foo"`
		Baz int      `header:"baz"`
		Qux bool     `header:"qux,default=true"`
	}
	r := httptest.NewRequest(http.MethodGet, "/any", http.NoBody)
	r.Header.Set("foo", "bar")
	r.Header.Add("foo", "bar1")
	r.Header.Set("baz", "1")
	assert.Nil(t, ParseHeaders(r.Header, &val))
	assert.Equal(t, []string{"bar", "bar1"}, val.Foo)
	assert.Equal(t, 1, val.Baz)
	assert.True(t, val.Qux)
}

func TestParseHeadersArrayInt(t *testing.T) {
	var val struct {
		Foo []int `header:"foo"`
	}
	r := httptest.NewRequest(http.MethodGet, "/any", http.NoBody)
	r.Header.Set("foo", "1")
	r.Header.Add("foo", "2")
	assert.Nil(t, ParseHeaders(r.Header, &val))
	assert.Equal(t, []int{1, 2}, val.Foo)
}
