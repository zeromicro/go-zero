package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorsHandlerWithOrigins(t *testing.T) {
	tests := []struct {
		name      string
		origins   []string
		reqOrigin string
		expect    string
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
			name:      "allow sub origins",
			origins:   []string{"local", "remote"},
			reqOrigin: "sub.local",
			expect:    "sub.local",
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
	}

	methods := []string{
		http.MethodOptions,
		http.MethodGet,
		http.MethodPost,
	}

	for _, test := range tests {
		for _, method := range methods {
			test := test
			t.Run(test.name+"-handler", func(t *testing.T) {
				r := httptest.NewRequest(method, "http://localhost", http.NoBody)
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
				r := httptest.NewRequest(method, "http://localhost", http.NoBody)
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
		}
	}

	for _, test := range tests {
		for _, method := range methods {
			test := test
			t.Run(test.name+"-middleware", func(t *testing.T) {
				r := httptest.NewRequest(method, "http://localhost", http.NoBody)
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
				r := httptest.NewRequest(method, "http://localhost", http.NoBody)
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
		}
	}
}
