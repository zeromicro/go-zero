package rest

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/rest/httpx"
	"github.com/tal-tech/go-zero/rest/router"
)

func TestNewServer(t *testing.T) {
	const configYaml = `
Name: foo
Port: 54321
`
	var cnf RestConf
	assert.Nil(t, conf.LoadConfigFromYamlBytes([]byte(configYaml), &cnf))
	failStart := func(server *Server) {
		server.opts.start = func(e *engine) error {
			return http.ErrServerClosed
		}
	}

	tests := []struct {
		c    RestConf
		opts []RunOption
		fail bool
	}{
		{
			c:    RestConf{},
			opts: []RunOption{failStart},
			fail: true,
		},
		{
			c:    cnf,
			opts: []RunOption{failStart},
		},
		{
			c:    cnf,
			opts: []RunOption{WithNotAllowedHandler(nil), failStart},
		},
		{
			c:    cnf,
			opts: []RunOption{WithNotFoundHandler(nil), failStart},
		},
		{
			c:    cnf,
			opts: []RunOption{WithUnauthorizedCallback(nil), failStart},
		},
		{
			c:    cnf,
			opts: []RunOption{WithUnsignedCallback(nil), failStart},
		},
	}

	for _, test := range tests {
		srv, err := NewServer(test.c, test.opts...)
		if test.fail {
			assert.NotNil(t, err)
		}
		if err != nil {
			continue
		}

		srv.Use(ToMiddleware(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}))
		srv.AddRoute(Route{
			Method:  http.MethodGet,
			Path:    "/",
			Handler: nil,
		}, WithJwt("thesecret"), WithSignature(SignatureConf{}),
			WithJwtTransition("preivous", "thenewone"))
		srv.Start()
		srv.Stop()
	}
}

func TestWithMiddleware(t *testing.T) {
	m := make(map[string]string)
	rt := router.NewRouter()
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
		assert.Nil(t, rt.Handle(route.Method, route.Path, route.Handler))
	}
	for _, url := range urls {
		r, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		rr := httptest.NewRecorder()
		rt.ServeHTTP(rr, r)

		assert.Equal(t, "whatever:200000", rr.Body.String())
	}

	assert.EqualValues(t, map[string]string{
		"kevin": "2017",
		"wan":   "2020",
	}, m)
}

func TestMultiMiddlewares(t *testing.T) {
	m := make(map[string]string)
	rt := router.NewRouter()
	handler := func(w http.ResponseWriter, r *http.Request) {
		var v struct {
			Nickname string `form:"nickname"`
			Zipcode  int64  `form:"zipcode"`
		}

		err := httpx.Parse(r, &v)
		assert.Nil(t, err)
		_, err = io.WriteString(w, fmt.Sprintf("%s:%s", v.Nickname, m[v.Nickname]))
		assert.Nil(t, err)
	}
	rs := WithMiddlewares([]Middleware{
		func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				var v struct {
					Name string `path:"name"`
					Year string `path:"year"`
				}
				assert.Nil(t, httpx.ParsePath(r, &v))
				m[v.Name] = v.Year
				next.ServeHTTP(w, r)
			}
		},
		func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				var v struct {
					Name    string `form:"nickname"`
					Zipcode string `form:"zipcode"`
				}
				assert.Nil(t, httpx.ParseForm(r, &v))
				assert.NotEmpty(t, m)
				m[v.Name] = v.Zipcode + v.Zipcode
				next.ServeHTTP(w, r)
			}
		},
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
		assert.Nil(t, rt.Handle(route.Method, route.Path, route.Handler))
	}
	for _, url := range urls {
		r, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err)

		rr := httptest.NewRecorder()
		rt.ServeHTTP(rr, r)

		assert.Equal(t, "whatever:200000200000", rr.Body.String())
	}

	assert.EqualValues(t, map[string]string{
		"kevin":    "2017",
		"wan":      "2020",
		"whatever": "200000200000",
	}, m)
}

func TestWithPriority(t *testing.T) {
	var fr featuredRoutes
	WithPriority()(&fr)
	assert.True(t, fr.priority)
}
