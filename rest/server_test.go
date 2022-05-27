package rest

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/router"
)

func TestNewServer(t *testing.T) {
	writer := logx.Reset()
	defer logx.SetWriter(writer)
	logx.SetWriter(logx.NewWriter(ioutil.Discard))

	const configYaml = `
Name: foo
Port: 54321
`
	var cnf RestConf
	assert.Nil(t, conf.LoadFromYamlBytes([]byte(configYaml), &cnf))

	tests := []struct {
		c    RestConf
		opts []RunOption
		fail bool
	}{
		{
			c:    RestConf{},
			opts: []RunOption{WithRouter(mockedRouter{}), WithCors()},
		},
		{
			c:    cnf,
			opts: []RunOption{WithRouter(mockedRouter{})},
		},
		{
			c:    cnf,
			opts: []RunOption{WithRouter(mockedRouter{}), WithNotAllowedHandler(nil)},
		},
		{
			c:    cnf,
			opts: []RunOption{WithNotFoundHandler(nil), WithRouter(mockedRouter{})},
		},
		{
			c:    cnf,
			opts: []RunOption{WithUnauthorizedCallback(nil), WithRouter(mockedRouter{})},
		},
		{
			c:    cnf,
			opts: []RunOption{WithUnsignedCallback(nil), WithRouter(mockedRouter{})},
		},
	}

	for _, test := range tests {
		var svr *Server
		var err error
		if test.fail {
			_, err = NewServer(test.c, test.opts...)
			assert.NotNil(t, err)
			continue
		} else {
			svr = MustNewServer(test.c, test.opts...)
		}

		svr.Use(ToMiddleware(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}))
		svr.AddRoute(Route{
			Method:  http.MethodGet,
			Path:    "/",
			Handler: nil,
		}, WithJwt("thesecret"), WithSignature(SignatureConf{}),
			WithJwtTransition("preivous", "thenewone"))

		func() {
			defer func() {
				p := recover()
				switch v := p.(type) {
				case error:
					assert.Equal(t, "foo", v.Error())
				default:
					t.Fail()
				}
			}()

			svr.Start()
			svr.Stop()
		}()
	}
}

func TestWithMaxBytes(t *testing.T) {
	const maxBytes = 1000
	var fr featuredRoutes
	WithMaxBytes(maxBytes)(&fr)
	assert.Equal(t, int64(maxBytes), fr.maxBytes)
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
		ToMiddleware(func(next http.Handler) http.Handler {
			return next
		}),
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

func TestWithPrefix(t *testing.T) {
	fr := featuredRoutes{
		routes: []Route{
			{
				Path: "/hello",
			},
			{
				Path: "/world",
			},
		},
	}
	WithPrefix("/api")(&fr)
	var vals []string
	for _, r := range fr.routes {
		vals = append(vals, r.Path)
	}
	assert.EqualValues(t, []string{"/api/hello", "/api/world"}, vals)
}

func TestWithPriority(t *testing.T) {
	var fr featuredRoutes
	WithPriority()(&fr)
	assert.True(t, fr.priority)
}

func TestWithTimeout(t *testing.T) {
	var fr featuredRoutes
	WithTimeout(time.Hour)(&fr)
	assert.Equal(t, time.Hour, fr.timeout)
}

func TestWithTLSConfig(t *testing.T) {
	const configYaml = `
Name: foo
Port: 54321
`
	var cnf RestConf
	assert.Nil(t, conf.LoadFromYamlBytes([]byte(configYaml), &cnf))

	testConfig := &tls.Config{
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	testCases := []struct {
		c    RestConf
		opts []RunOption
		res  *tls.Config
	}{
		{
			c:    cnf,
			opts: []RunOption{WithTLSConfig(testConfig)},
			res:  testConfig,
		},
		{
			c:    cnf,
			opts: []RunOption{WithUnsignedCallback(nil)},
			res:  nil,
		},
	}

	for _, testCase := range testCases {
		svr, err := NewServer(testCase.c, testCase.opts...)
		assert.Nil(t, err)
		assert.Equal(t, svr.ngin.tlsConfig, testCase.res)
	}
}

func TestWithCors(t *testing.T) {
	const configYaml = `
Name: foo
Port: 54321
`
	var cnf RestConf
	assert.Nil(t, conf.LoadFromYamlBytes([]byte(configYaml), &cnf))
	rt := router.NewRouter()
	svr, err := NewServer(cnf, WithRouter(rt))
	assert.Nil(t, err)

	opt := WithCors("local")
	opt(svr)
}

func TestWithCustomCors(t *testing.T) {
	const configYaml = `
Name: foo
Port: 54321
`
	var cnf RestConf
	assert.Nil(t, conf.LoadFromYamlBytes([]byte(configYaml), &cnf))
	rt := router.NewRouter()
	svr, err := NewServer(cnf, WithRouter(rt))
	assert.Nil(t, err)

	opt := WithCustomCors(func(header http.Header) {
		header.Set("foo", "bar")
	}, func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusOK)
	}, "local")
	opt(svr)
}
