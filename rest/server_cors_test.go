package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/logx/logtest"
	"github.com/zeromicro/go-zero/rest/router"
)

func TestCorsOptionsPreserveOriginPolicy(t *testing.T) {
	logtest.Discard(t)
	const trusted = "https://trusted.example"
	custom := WithCustomCors(func(h http.Header) {
		h.Set("X-Policy", "custom")
	}, func(w http.ResponseWriter) {
		w.Header().Set("X-Not-Allowed", "custom")
	}, trusted)

	for _, test := range []struct {
		name   string
		opts   []RunOption
		custom bool
	}{
		{name: "origins then headers", opts: []RunOption{WithCors(trusted), WithCorsHeaders("X-Custom")}},
		{name: "headers then origins", opts: []RunOption{WithCorsHeaders("X-Custom"), WithCors(trusted)}},
		{name: "custom then headers", opts: []RunOption{custom, WithCorsHeaders("X-Custom")}, custom: true},
		{name: "headers then custom", opts: []RunOption{WithCorsHeaders("X-Custom"), custom}, custom: true},
		{name: "repeated origin policies", opts: []RunOption{WithCors("https://old.example"), WithCorsHeaders("X-Custom"), WithCors(trusted)}},
		{name: "file router between options", opts: []RunOption{WithCors(trusted), WithFileServer("/assets", http.Dir("testdata")), WithCorsHeaders("X-Custom")}},
	} {
		t.Run(test.name, func(t *testing.T) {
			s, err := NewServer(RestConf{}, test.opts...)
			require.NoError(t, err)
			s.AddRoute(Route{Method: http.MethodGet, Path: "/private", Handler: func(w http.ResponseWriter, r *http.Request) {
				cookie, err := r.Cookie("session")
				require.NoError(t, err)
				require.Equal(t, "alice", cookie.Value)
				_, _ = w.Write([]byte("private-for-alice"))
			}})
			require.NoError(t, s.build())

			for _, origin := range []string{trusted, "https://attacker.example", "https://old.example", "null"} {
				for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodOptions} {
					for _, path := range []string{"/private", "/missing"} {
						t.Run(origin+"/"+method+path, func(t *testing.T) {
							r := httptest.NewRequest(method, path, nil)
							r.Header.Set("Origin", origin)
							r.Header.Set("Access-Control-Request-Method", http.MethodGet)
							r.Header.Set("Access-Control-Request-Headers", "X-Custom")
							r.AddCookie(&http.Cookie{Name: "session", Value: "alice"})
							w := httptest.NewRecorder()
							s.serve(w, r)
							if origin == trusted {
								assert.Equal(t, trusted, w.Header().Get("Access-Control-Allow-Origin"))
								assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
								assert.Contains(t, strings.Join(w.Header().Values("Access-Control-Allow-Headers"), ","), "X-Custom")
							} else {
								assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
								assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"))
							}
							if test.custom {
								assert.Equal(t, "custom", w.Header().Get("X-Policy"))
								if method == http.MethodPost && path == "/private" {
									assert.Equal(t, "custom", w.Header().Get("X-Not-Allowed"))
								}
							}
							if method == http.MethodOptions {
								assert.Equal(t, http.StatusNoContent, w.Code)
							} else if method == http.MethodGet && path == "/private" {
								assert.Equal(t, "private-for-alice", w.Body.String())
							}
						})
					}
				}
			}
		})
	}
}

func TestCorsRoutingOptions(t *testing.T) {
	logtest.Discard(t)
	const trusted = "https://trusted.example"
	assets := http.FS(fstest.MapFS{"hello.txt": &fstest.MapFile{Data: []byte("hello")}})
	for _, test := range []struct {
		name string
		opts []RunOption
	}{
		{name: "file server before cors", opts: []RunOption{WithFileServer("/assets", assets), WithCors(trusted), WithCorsHeaders("X-Custom")}},
		{name: "file server between cors options", opts: []RunOption{WithCors(trusted), WithFileServer("/assets", assets), WithCorsHeaders("X-Custom")}},
		{name: "file server after cors", opts: []RunOption{WithCors(trusted), WithCorsHeaders("X-Custom"), WithFileServer("/assets", assets)}},
		{name: "router between cors options", opts: []RunOption{WithCors(trusted), WithRouter(router.NewRouter()), WithCorsHeaders("X-Custom"), WithFileServer("/assets", assets)}},
		{name: "router after cors", opts: []RunOption{WithCors(trusted), WithCorsHeaders("X-Custom"), WithRouter(router.NewRouter()), WithFileServer("/assets", assets)}},
		{name: "router after headers before origins", opts: []RunOption{WithCorsHeaders("X-Custom"), WithRouter(router.NewRouter()), WithCors(trusted), WithFileServer("/assets", assets)}},
		{name: "multiple file servers", opts: []RunOption{WithCors(trusted), WithFileServer("/assets", assets), WithCorsHeaders("X-Custom"), WithFileServer("/other", assets)}},
	} {
		t.Run(test.name, func(t *testing.T) {
			s, err := NewServer(RestConf{}, test.opts...)
			require.NoError(t, err)
			s.AddRoute(Route{Method: http.MethodGet, Path: "/private", Handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("private"))
			}})
			require.NoError(t, s.build())
			paths := []string{"/private", "/assets/hello.txt"}
			if test.name == "multiple file servers" {
				paths = append(paths, "/other/hello.txt")
			}
			for _, path := range paths {
				for _, method := range []string{http.MethodGet, http.MethodOptions, http.MethodPost} {
					for _, origin := range []string{trusted, "https://attacker.example"} {
						t.Run(method+path+"/"+origin, func(t *testing.T) {
							r := httptest.NewRequest(method, path, nil)
							r.Header.Set("Origin", origin)
							r.Header.Set("Access-Control-Request-Method", http.MethodGet)
							r.Header.Set("Access-Control-Request-Headers", "X-Custom")
							w := httptest.NewRecorder()
							s.serve(w, r)
							if origin == trusted {
								assert.Equal(t, trusted, w.Header().Get("Access-Control-Allow-Origin"))
								assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
								assert.Contains(t, strings.Join(w.Header().Values("Access-Control-Allow-Headers"), ","), "X-Custom")
							} else {
								assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
								assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"))
							}
							if method == http.MethodGet {
								assert.Equal(t, http.StatusOK, w.Code)
								if path == "/private" {
									assert.Equal(t, "private", w.Body.String())
								} else {
									assert.Equal(t, "hello", w.Body.String())
								}
							} else if method == http.MethodOptions {
								assert.Equal(t, http.StatusNoContent, w.Code)
							}
						})
					}
				}
			}
		})
	}
}

func TestCorsWildcardDoesNotAllowCredentials(t *testing.T) {
	logtest.Discard(t)
	for _, opts := range [][]RunOption{
		{WithCors()},
		{WithCors("*")},
		{WithCors("https://trusted.example", "*")},
		{WithCorsHeaders("X-Custom"), WithCorsHeaders("X-Another")},
	} {
		s, err := NewServer(RestConf{}, opts...)
		require.NoError(t, err)
		for _, method := range []string{http.MethodGet, http.MethodOptions} {
			r := httptest.NewRequest(method, "/missing", nil)
			r.Header.Set("Origin", "https://attacker.example")
			w := httptest.NewRecorder()
			s.serve(w, r)
			assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
			assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"))
			if len(opts) == 2 {
				headers := strings.Join(w.Header().Values("Access-Control-Allow-Headers"), ",")
				assert.Contains(t, headers, "X-Custom")
				assert.Contains(t, headers, "X-Another")
			}
		}
	}
}

func TestWithCorsKeepsCustomCallbacks(t *testing.T) {
	logtest.Discard(t)
	const trusted = "https://trusted.example"
	s, err := NewServer(RestConf{},
		WithCustomCors(func(h http.Header) { h.Set("X-Policy", "custom") },
			func(w http.ResponseWriter) { w.Header().Set("X-Not-Allowed", "custom") }, "https://old.example"),
		WithCors(trusted))
	require.NoError(t, err)
	s.AddRoute(Route{Method: http.MethodGet, Path: "/ok", Handler: func(w http.ResponseWriter, r *http.Request) {}})
	require.NoError(t, s.build())

	r := httptest.NewRequest(http.MethodGet, "/ok", nil)
	r.Header.Set("Origin", trusted)
	w := httptest.NewRecorder()
	s.serve(w, r)
	assert.Equal(t, trusted, w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "custom", w.Header().Get("X-Policy"))

	r = httptest.NewRequest(http.MethodPost, "/ok", nil)
	r.Header.Set("Origin", "https://old.example")
	w = httptest.NewRecorder()
	s.serve(w, r)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "custom", w.Header().Get("X-Not-Allowed"))
}
