package rest

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/rest/httpx"
	"github.com/tal-tech/go-zero/rest/router"
)

func TestWithMiddleware(t *testing.T) {
	m := make(map[string]string)
	router := router.NewPatRouter()
	handler := func(w http.ResponseWriter, r *http.Request) {
		var v struct {
			Nickname string `form:"nickname"`
			Zipcode  int64  `form:"zipcode"`
		}

		err := httpx.Parse(r, &v)
		assert.Nil(t, err)
		_, err = io.WriteString(w, fmt.Sprintf("%s:%d", v.Nickname, v.Zipcode))
		assert.Nil(t, err)
	}
	rs := WithMiddleware(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var v struct {
				Name string `path:"name"`
				Year string `path:"year"`
			}
			assert.Nil(t, httpx.ParsePath(r, &v))
			m[v.Name] = v.Year
			next.ServeHTTP(w, r)
		}
	}, Route{
		Method:  http.MethodGet,
		Path:    "/first/:name/:year",
		Handler: handler,
	}, Route{
		Method:  http.MethodGet,
		Path:    "/second/:name/:year",
		Handler: handler,
	})

	urls := []string{
		"http://hello.com/first/kevin/2017?nickname=whatever&zipcode=200000",
		"http://hello.com/second/wan/2020?nickname=whatever&zipcode=200000",
	}
	for _, route := range rs {
		assert.Nil(t, router.Handle(route.Method, route.Path, route.Handler))
	}
	for _, url := range urls {
		r, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, r)

		assert.Equal(t, "whatever:200000", rr.Body.String())
	}

	assert.EqualValues(t, map[string]string{
		"kevin": "2017",
		"wan":   "2020",
	}, m)
}
