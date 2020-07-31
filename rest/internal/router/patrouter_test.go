package router

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"zero/rest/internal/context"
)

type mockedResponseWriter struct {
	code int
}

func (m *mockedResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m *mockedResponseWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (m *mockedResponseWriter) WriteHeader(code int) {
	m.code = code
}

func TestPatRouterHandleErrors(t *testing.T) {
	tests := []struct {
		method string
		path   string
		err    error
	}{
		{"FAKE", "", ErrInvalidMethod},
		{"GET", "", ErrInvalidPath},
	}

	for _, test := range tests {
		t.Run(test.method, func(t *testing.T) {
			router := NewPatRouter()
			err := router.Handle(test.method, test.path, nil)
			assert.Error(t, ErrInvalidMethod, err)
		})
	}
}

func TestPatRouterNotFound(t *testing.T) {
	var notFound bool
	router := NewPatRouter()
	router.SetNotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notFound = true
	}))
	router.Handle(http.MethodGet, "/a/b", nil)
	r, _ := http.NewRequest(http.MethodGet, "/b/c", nil)
	w := new(mockedResponseWriter)
	router.ServeHTTP(w, r)
	assert.True(t, notFound)
}

func TestPatRouter(t *testing.T) {
	tests := []struct {
		method string
		path   string
		expect bool
		code   int
		err    error
	}{
		// we don't explicitly set status code, framework will do it.
		{http.MethodGet, "/a/b", true, 0, nil},
		{http.MethodGet, "/a/b/", true, 0, nil},
		{http.MethodGet, "/a/b?a=b", true, 0, nil},
		{http.MethodGet, "/a/b/?a=b", true, 0, nil},
		{http.MethodGet, "/a/b/c?a=b", true, 0, nil},
		{http.MethodGet, "/b/d", false, http.StatusNotFound, nil},
	}

	for _, test := range tests {
		t.Run(test.method+":"+test.path, func(t *testing.T) {
			routed := false
			router := NewPatRouter()
			err := router.Handle(test.method, "/a/:b", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				routed = true
				assert.Equal(t, 1, len(context.Vars(r)))
			}))
			assert.Nil(t, err)
			err = router.Handle(test.method, "/a/b/c", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				routed = true
				assert.Nil(t, context.Vars(r))
			}))
			assert.Nil(t, err)
			err = router.Handle(test.method, "/b/c", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				routed = true
			}))
			assert.Nil(t, err)

			w := new(mockedResponseWriter)
			r, _ := http.NewRequest(test.method, test.path, nil)
			router.ServeHTTP(w, r)
			assert.Equal(t, test.expect, routed)
			assert.Equal(t, test.code, w.code)

			if test.code == 0 {
				r, _ = http.NewRequest(http.MethodPut, test.path, nil)
				router.ServeHTTP(w, r)
				assert.Equal(t, http.StatusMethodNotAllowed, w.code)
			}
		})
	}
}

func BenchmarkPatRouter(b *testing.B) {
	b.ReportAllocs()

	router := NewPatRouter()
	router.Handle(http.MethodGet, "/api/:user/:name", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	w := &mockedResponseWriter{}
	r, _ := http.NewRequest(http.MethodGet, "/api/a/b", nil)
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, r)
	}
}
