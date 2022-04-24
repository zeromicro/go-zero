package cors

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorsHandlerWithOrigins(t *testing.T) {
	tests := []struct {
		name       string
		origins    []string
		reqOrigin  string
		expect     string
		reqOrigins []string
	}{
		{
			name:   "allow all origins",
			expect: allOrigins,
		},
		{
			name:      "allow one origin",
			origins:   []string{"http://local"},
			reqOrigin: "http://local",
			expect:    "http://local",
		},
		{
			name:      "allow many origins",
			origins:   []string{"http://local", "http://remote"},
			reqOrigin: "http://local",
			expect:    "http://local",
		},
		{
			name:      "allow all origins",
			reqOrigin: "http://local",
			expect:    "*",
		},
		{
			name:      "allow many origins with all mark",
			origins:   []string{"http://local", "http://remote", "*"},
			reqOrigin: "http://another",
			expect:    "http://another",
		},
		{
			name:      "not allow origin",
			origins:   []string{"http://local", "http://remote"},
			reqOrigin: "http://another",
		},
		{
			name:    "allow regexp origin",
			origins: []string{"127\\.0\\.0\\.1((:(\\d+))*)$", "localhost((:(\\d+))*)$", "(.+)\\.go-zero\\.com$"},
			reqOrigins: []string{
				"http://127.0.0.1",
				"http://127.0.0.1:3000",
				"http://localhost",
				"http://localhost:8000",
				"https://test.go-zero.com",
				"https://three.test.go-zero.com",
			},
		},
		{
			name:    "not allow regexp origin",
			origins: []string{"127\\.0\\.0\\.1((:(\\d+))*)$", "localhost((:(\\d+))*)$", "(.+)\\.go-zero\\.com$"},
			reqOrigins: []string{
				"http://128.0.0.1",
				"http://128.0.0.1:3000",
				"http://local",
				"http://local:8000",
				"https://test.go-zero.com.test.com",
				"https://three.test.go-zero.com.baidu.com",
			},
		},
	}

	methods := []string{
		http.MethodOptions,
		http.MethodGet,
		http.MethodPost,
	}

	for _, test := range tests {
		for _, method := range methods {
			test := test
			if len(test.reqOrigins) > 0 {
				t.Run(test.name+"-handler", func(t *testing.T) {
					r := httptest.NewRequest(method, "http://localhost", nil)
					r.Header.Set(originHeader, test.reqOrigin)
					w := httptest.NewRecorder()
					handler := NotAllowedHandler(nil, test.origins...)
					handler.ServeHTTP(w, r)
					if method == http.MethodOptions {
						assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
					} else {
						assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
					}
					assert.Equal(t, test.expect, w.Header().Get(allowOrigin))
				})

				t.Run(test.name+"-handler-custom", func(t *testing.T) {
					r := httptest.NewRequest(method, "http://localhost", nil)
					r.Header.Set(originHeader, test.reqOrigin)
					w := httptest.NewRecorder()
					handler := NotAllowedHandler(func(w http.ResponseWriter) {
						w.Header().Set("foo", "bar")
					}, test.origins...)
					handler.ServeHTTP(w, r)
					if method == http.MethodOptions {
						assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
					} else {
						assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
					}
					assert.Equal(t, test.expect, w.Header().Get(allowOrigin))
					assert.Equal(t, "bar", w.Header().Get("foo"))
				})
			} else {
				t.Run(test.name+"-handler", func(t *testing.T) {
					for _, reqOrigin := range test.reqOrigins {
						r := httptest.NewRequest(method, "http://localhost", nil)
						r.Header.Set(originHeader, reqOrigin)
						w := httptest.NewRecorder()
						handler := NotAllowedHandler(nil, test.origins...)
						handler.ServeHTTP(w, r)
						if method == http.MethodOptions {
							assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
						} else {
							assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
						}
						var matched bool
						for _, expression := range test.origins {
							reg := regexp.MustCompile(expression)
							if reg.MatchString(w.Header().Get(allowOrigin)) {
								matched = true
								break
							}
						}
						assert.Equal(t, true, matched, w.Header().Get(allowOrigin))
					}
				})

				t.Run(test.name+"-handler-custom", func(t *testing.T) {
					for _, reqOrigin := range test.reqOrigins {
						r := httptest.NewRequest(method, "http://localhost", nil)
						r.Header.Set(originHeader, reqOrigin)
						w := httptest.NewRecorder()
						handler := NotAllowedHandler(func(w http.ResponseWriter) {
							w.Header().Set("foo", "bar")
						}, test.origins...)
						handler.ServeHTTP(w, r)
						if method == http.MethodOptions {
							assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
						} else {
							assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
						}
						var matched bool
						for _, expression := range test.origins {
							reg := regexp.MustCompile(expression)
							if reg.MatchString(w.Header().Get(allowOrigin)) {
								matched = true
								break
							}
						}
						assert.Equal(t, false, matched)
						assert.Equal(t, "bar", w.Header().Get("foo"))
					}
				})
			}
		}
	}

	for _, test := range tests {
		for _, method := range methods {
			test := test
			if len(test.reqOrigins) > 0 {
				t.Run(test.name+"-middleware", func(t *testing.T) {
					r := httptest.NewRequest(method, "http://localhost", nil)
					r.Header.Set(originHeader, test.reqOrigin)
					w := httptest.NewRecorder()
					handler := Middleware(nil, test.origins...)(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
					})
					handler.ServeHTTP(w, r)
					if method == http.MethodOptions {
						assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
					} else {
						assert.Equal(t, http.StatusOK, w.Result().StatusCode)
					}
					assert.Equal(t, test.expect, w.Header().Get(allowOrigin))
				})
				t.Run(test.name+"-middleware-custom", func(t *testing.T) {
					r := httptest.NewRequest(method, "http://localhost", nil)
					r.Header.Set(originHeader, test.reqOrigin)
					w := httptest.NewRecorder()
					handler := Middleware(func(header http.Header) {
						header.Set("foo", "bar")
					}, test.origins...)(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
					})
					handler.ServeHTTP(w, r)
					if method == http.MethodOptions {
						assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
					} else {
						assert.Equal(t, http.StatusOK, w.Result().StatusCode)
					}
					assert.Equal(t, test.expect, w.Header().Get(allowOrigin))
					assert.Equal(t, "bar", w.Header().Get("foo"))
				})
			} else {
				t.Run(test.name+"-middleware", func(t *testing.T) {
					for _, reqOrigin := range test.reqOrigins {
						r := httptest.NewRequest(method, "http://localhost", nil)
						r.Header.Set(originHeader, reqOrigin)
						w := httptest.NewRecorder()
						handler := Middleware(nil, test.origins...)(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusOK)
						})
						handler.ServeHTTP(w, r)
						if method == http.MethodOptions {
							assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
						} else {
							assert.Equal(t, http.StatusOK, w.Result().StatusCode)
						}
						var matched bool
						for _, expression := range test.origins {
							reg := regexp.MustCompile(expression)
							if reg.MatchString(w.Header().Get(allowOrigin)) {
								matched = true
								break
							}
						}
						assert.Equal(t, true, matched, w.Header().Get(allowOrigin))
					}
				})
				t.Run(test.name+"-middleware-custom", func(t *testing.T) {
					for _, reqOrigin := range test.reqOrigins {
						r := httptest.NewRequest(method, "http://localhost", nil)
						r.Header.Set(originHeader, reqOrigin)
						w := httptest.NewRecorder()
						handler := Middleware(func(header http.Header) {
							header.Set("foo", "bar")
						}, test.origins...)(func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusOK)
						})
						handler.ServeHTTP(w, r)
						if method == http.MethodOptions {
							assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
						} else {
							assert.Equal(t, http.StatusOK, w.Result().StatusCode)
						}
						var matched bool
						for _, expression := range test.origins {
							reg := regexp.MustCompile(expression)
							if reg.MatchString(w.Header().Get(allowOrigin)) {
								matched = true
								break
							}
						}
						assert.Equal(t, false, matched)
						assert.Equal(t, "bar", w.Header().Get("foo"))
					}
				})
			}
		}
	}
}
