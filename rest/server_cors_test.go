package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/logx/logtest"
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
