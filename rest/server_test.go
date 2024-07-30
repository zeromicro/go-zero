package rest

import (
	"crypto/tls"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx/logtest"
	"github.com/zeromicro/go-zero/rest/chain"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/internal/cors"
	"github.com/zeromicro/go-zero/rest/router"
)

const (
	exampleContent = "example content"
	sampleContent  = "sample content"
)

func TestNewServer(t *testing.T) {
	logtest.Discard(t)

	const configYaml = `
Name: foo
Host: localhost
Port: 0
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

			svr.StartWithOpts(func(svr *http.Server) {
				svr.RegisterOnShutdown(func() {})
			})
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

func TestWithFileServerMiddleware(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		dir             string
		requestPath     string
		expectedStatus  int
		expectedContent string
	}{
		{
			name:            "Serve static file",
			path:            "/assets/",
			dir:             "./testdata",
			requestPath:     "/assets/example.txt",
			expectedStatus:  http.StatusOK,
			expectedContent: exampleContent,
		},
		{
			name:           "Pass through non-matching path",
			path:           "/static/",
			dir:            "./testdata",
			requestPath:    "/other/path",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:            "Directory with trailing slash",
			path:            "/static",
			dir:             "testdata",
			requestPath:     "/static/sample.txt",
			expectedStatus:  http.StatusOK,
			expectedContent: sampleContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := MustNewServer(RestConf{}, WithFileServer(tt.path, http.Dir(tt.dir)))
			req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
			rr := httptest.NewRecorder()

			server.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if len(tt.expectedContent) > 0 {
				assert.Equal(t, tt.expectedContent, rr.Body.String())
			}
		})
	}
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
	vals := make([]string, 0, len(fr.routes))
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
	defer svr.Stop()

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

func TestWithCorsHeaders(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
	}{
		{
			name:    "single header",
			headers: []string{"UserHeader"},
		},
		{
			name:    "multiple headers",
			headers: []string{"UserHeader", "X-Requested-With"},
		},
		{
			name:    "no headers",
			headers: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const configYaml = `
Name: foo
Port: 54321
`
			var cnf RestConf
			assert.Nil(t, conf.LoadFromYamlBytes([]byte(configYaml), &cnf))
			rt := router.NewRouter()
			svr, err := NewServer(cnf, WithRouter(rt))
			assert.Nil(t, err)
			defer svr.Stop()
			option := WithCorsHeaders(tt.headers...)
			option(svr)

			// Assuming newCorsRouter sets headers correctly,
			// we would need to verify the behavior here. Since we don't have
			// direct access to headers, we'll mock newCorsRouter to capture it.
			w := httptest.NewRecorder()
			svr.ServeHTTP(w, httptest.NewRequest(http.MethodOptions, "/", nil))

			vals := w.Header().Values("Access-Control-Allow-Headers")
			respHeaders := make(map[string]struct{})
			for _, header := range vals {
				headers := strings.Split(header, ", ")
				for _, h := range headers {
					if len(h) > 0 {
						respHeaders[h] = struct{}{}
					}
				}
			}
			for _, h := range tt.headers {
				_, ok := respHeaders[h]
				assert.Truef(t, ok, "expected header %s not found", h)
			}
		})
	}
}

func TestServer_PrintRoutes(t *testing.T) {
	const (
		configYaml = `
Name: foo
Port: 54321
`
		expect = `Routes:
  GET /bar
  GET /foo
  GET /foo/:bar
  GET /foo/:bar/baz
`
	)

	var cnf RestConf
	assert.Nil(t, conf.LoadFromYamlBytes([]byte(configYaml), &cnf))

	svr, err := NewServer(cnf)
	assert.Nil(t, err)

	svr.AddRoutes([]Route{
		{
			Method:  http.MethodGet,
			Path:    "/foo",
			Handler: http.NotFound,
		},
		{
			Method:  http.MethodGet,
			Path:    "/bar",
			Handler: http.NotFound,
		},
		{
			Method:  http.MethodGet,
			Path:    "/foo/:bar",
			Handler: http.NotFound,
		},
		{
			Method:  http.MethodGet,
			Path:    "/foo/:bar/baz",
			Handler: http.NotFound,
		},
	})

	old := os.Stdout
	r, w, err := os.Pipe()
	assert.Nil(t, err)
	os.Stdout = w
	defer func() {
		os.Stdout = old
	}()

	svr.PrintRoutes()
	ch := make(chan string)

	go func() {
		var buf strings.Builder
		io.Copy(&buf, r)
		ch <- buf.String()
	}()

	w.Close()
	out := <-ch
	assert.Equal(t, expect, out)
}

func TestServer_Routes(t *testing.T) {
	const (
		configYaml = `
Name: foo
Port: 54321
`
		expect = `GET /foo GET /bar GET /foo/:bar GET /foo/:bar/baz`
	)

	var cnf RestConf
	assert.Nil(t, conf.LoadFromYamlBytes([]byte(configYaml), &cnf))

	svr, err := NewServer(cnf)
	assert.Nil(t, err)

	svr.AddRoutes([]Route{
		{
			Method:  http.MethodGet,
			Path:    "/foo",
			Handler: http.NotFound,
		},
		{
			Method:  http.MethodGet,
			Path:    "/bar",
			Handler: http.NotFound,
		},
		{
			Method:  http.MethodGet,
			Path:    "/foo/:bar",
			Handler: http.NotFound,
		},
		{
			Method:  http.MethodGet,
			Path:    "/foo/:bar/baz",
			Handler: http.NotFound,
		},
	})

	routes := svr.Routes()
	var buf strings.Builder
	for i := 0; i < len(routes); i++ {
		buf.WriteString(routes[i].Method)
		buf.WriteString(" ")
		buf.WriteString(routes[i].Path)
		buf.WriteString(" ")
	}

	assert.Equal(t, expect, strings.Trim(buf.String(), " "))
}

func TestHandleError(t *testing.T) {
	assert.NotPanics(t, func() {
		handleError(nil)
		handleError(http.ErrServerClosed)
	})
}

func TestValidateSecret(t *testing.T) {
	assert.Panics(t, func() {
		validateSecret("short")
	})
}

func TestServer_WithChain(t *testing.T) {
	var called int32
	middleware1 := func() func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&called, 1)
				next.ServeHTTP(w, r)
				atomic.AddInt32(&called, 1)
			})
		}
	}
	middleware2 := func() func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&called, 1)
				next.ServeHTTP(w, r)
				atomic.AddInt32(&called, 1)
			})
		}
	}

	server := MustNewServer(RestConf{}, WithChain(chain.New(middleware1(), middleware2())))
	server.AddRoutes(
		[]Route{
			{
				Method: http.MethodGet,
				Path:   "/",
				Handler: func(_ http.ResponseWriter, _ *http.Request) {
					atomic.AddInt32(&called, 1)
				},
			},
		},
	)
	rt := router.NewRouter()
	assert.Nil(t, server.ngin.bindRoutes(rt))
	req, err := http.NewRequest(http.MethodGet, "/", http.NoBody)
	assert.Nil(t, err)
	rt.ServeHTTP(httptest.NewRecorder(), req)
	assert.Equal(t, int32(5), atomic.LoadInt32(&called))
}

func TestServer_WithCors(t *testing.T) {
	var called int32
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&called, 1)
			next.ServeHTTP(w, r)
		})
	}
	r := router.NewRouter()
	assert.Nil(t, r.Handle(http.MethodOptions, "/", middleware(http.NotFoundHandler())))

	cr := &corsRouter{
		Router:     r,
		middleware: cors.Middleware(nil, "*"),
	}
	req := httptest.NewRequest(http.MethodOptions, "/", http.NoBody)
	cr.ServeHTTP(httptest.NewRecorder(), req)
	assert.Equal(t, int32(0), atomic.LoadInt32(&called))
}

func TestServer_ServeHTTP(t *testing.T) {
	const configYaml = `
Name: foo
Port: 54321
`

	var cnf RestConf
	assert.Nil(t, conf.LoadFromYamlBytes([]byte(configYaml), &cnf))

	svr, err := NewServer(cnf)
	assert.Nil(t, err)

	svr.AddRoutes([]Route{
		{
			Method: http.MethodGet,
			Path:   "/foo",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				_, _ = writer.Write([]byte("succeed"))
				writer.WriteHeader(http.StatusOK)
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/bar",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				_, _ = writer.Write([]byte("succeed"))
				writer.WriteHeader(http.StatusOK)
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/user/:name",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				var userInfo struct {
					Name string `path:"name"`
				}

				err := httpx.Parse(request, &userInfo)
				if err != nil {
					_, _ = writer.Write([]byte("failed"))
					writer.WriteHeader(http.StatusBadRequest)
					return
				}

				_, _ = writer.Write([]byte("succeed"))
				writer.WriteHeader(http.StatusOK)
			},
		},
	})

	testCase := []struct {
		name string
		path string
		code int
	}{
		{
			name: "URI : /foo",
			path: "/foo",
			code: http.StatusOK,
		},
		{
			name: "URI : /bar",
			path: "/bar",
			code: http.StatusOK,
		},
		{
			name: "URI : undefined path",
			path: "/test",
			code: http.StatusNotFound,
		},
		{
			name: "URI : /user/:name",
			path: "/user/abc",
			code: http.StatusOK,
		},
	}

	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", test.path, nil)
			svr.ServeHTTP(w, req)
			assert.Equal(t, test.code, w.Code)
		})
	}
}

//go:embed testdata
var content embed.FS

func TestServerEmbedFileSystem(t *testing.T) {
	filesys, err := fs.Sub(content, "testdata")
	assert.NoError(t, err)

	server := MustNewServer(RestConf{}, WithFileServer("/assets", http.FS(filesys)))
	req, err := http.NewRequest(http.MethodGet, "/assets/sample.txt", http.NoBody)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)
	assert.Equal(t, sampleContent, rr.Body.String())
}
