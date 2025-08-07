package fileserver

import (
	"embed"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
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
			path:            "/static/",
			dir:             "./testdata",
			requestPath:     "/static/example.txt",
			expectedStatus:  http.StatusOK,
			expectedContent: "1",
		},
		{
			name:           "Pass through non-matching path",
			path:           "/static/",
			dir:            "./testdata",
			requestPath:    "/other/path",
			expectedStatus: http.StatusAlreadyReported,
		},
		{
			name:            "Directory with trailing slash",
			path:            "/assets",
			dir:             "testdata",
			requestPath:     "/assets/sample.txt",
			expectedStatus:  http.StatusOK,
			expectedContent: "2",
		},
		{
			name:           "Not exist file",
			path:           "/assets",
			dir:            "testdata",
			requestPath:    "/assets/not-exist.txt",
			expectedStatus: http.StatusAlreadyReported,
		},
		{
			name:           "Not exist file in root",
			path:           "/",
			dir:            "testdata",
			requestPath:    "/not-exist.txt",
			expectedStatus: http.StatusAlreadyReported,
		},
		{
			name:           "websocket request",
			path:           "/",
			dir:            "testdata",
			requestPath:    "/ws",
			expectedStatus: http.StatusAlreadyReported,
		},

		// http.FileServer redirects any request ending in "/index.html"
		// to the same path, without the final "index.html".
		{
			name:           "Serve index.html",
			path:           "/static",
			dir:            "testdata",
			requestPath:    "/static/index.html",
			expectedStatus: http.StatusMovedPermanently,
		},
		{
			name:           "Serve index.html with path with trailing slash",
			path:           "/static/",
			dir:            "testdata",
			requestPath:    "/static/index.html",
			expectedStatus: http.StatusMovedPermanently,
		},
		{
			name:           "Serve index.html in a nested directory",
			path:           "/static",
			dir:            "testdata",
			requestPath:    "/static/nested/index.html",
			expectedStatus: http.StatusMovedPermanently,
		},
		{
			name:            "Request index.html indirectly",
			path:            "/static",
			dir:             "testdata",
			requestPath:     "/static/",
			expectedStatus:  http.StatusOK,
			expectedContent: "hello",
		},
		{
			name:            "Request index.html in a nested directory indirectly",
			path:            "/static",
			dir:             "testdata",
			requestPath:     "/static/nested/",
			expectedStatus:  http.StatusOK,
			expectedContent: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := Middleware(tt.path, http.Dir(tt.dir))
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAlreadyReported)
			})

			handlerToTest := middleware(nextHandler)

			for i := 0; i < 2; i++ {
				req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
				rr := httptest.NewRecorder()

				handlerToTest.ServeHTTP(rr, req)

				assert.Equal(t, tt.expectedStatus, rr.Code)
				if len(tt.expectedContent) > 0 {
					assert.Equal(t, tt.expectedContent, rr.Body.String())
				}
			}
		})
	}
}

var (
	//go:embed testdata
	testdataFS embed.FS
)

func TestMiddleware_embedFS(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		requestPath     string
		expectedStatus  int
		expectedContent string
	}{
		{
			name:            "Serve static file",
			path:            "/static",
			requestPath:     "/static/example.txt",
			expectedStatus:  http.StatusOK,
			expectedContent: "1",
		},
		{
			name:            "Path with trailing slash",
			path:            "/static/",
			requestPath:     "/static/example.txt",
			expectedStatus:  http.StatusOK,
			expectedContent: "1",
		},
		{
			name:            "Root path",
			path:            "/",
			requestPath:     "/example.txt",
			expectedStatus:  http.StatusOK,
			expectedContent: "1",
		},
		{
			name:           "Pass through non-matching path",
			path:           "/static/",
			requestPath:    "/other/path",
			expectedStatus: http.StatusAlreadyReported,
		},
		{
			name:           "Not exist file",
			path:           "/assets",
			requestPath:    "/assets/not-exist.txt",
			expectedStatus: http.StatusAlreadyReported,
		},
		{
			name:           "Not exist file in root",
			path:           "/",
			requestPath:    "/not-exist.txt",
			expectedStatus: http.StatusAlreadyReported,
		},
		{
			name:           "websocket request",
			path:           "/",
			requestPath:    "/ws",
			expectedStatus: http.StatusAlreadyReported,
		},

		// http.FileServer redirects any request ending in "/index.html"
		// to the same path, without the final "index.html".
		{
			name:           "Serve index.html",
			path:           "/static",
			requestPath:    "/static/index.html",
			expectedStatus: http.StatusMovedPermanently,
		},
		{
			name:           "Serve index.html with path with trailing slash",
			path:           "/static/",
			requestPath:    "/static/index.html",
			expectedStatus: http.StatusMovedPermanently,
		},
		{
			name:           "Serve index.html in a nested directory",
			path:           "/static",
			requestPath:    "/static/nested/index.html",
			expectedStatus: http.StatusMovedPermanently,
		},
		{
			name:            "Request index.html indirectly",
			path:            "/static",
			requestPath:     "/static/",
			expectedStatus:  http.StatusOK,
			expectedContent: "hello",
		},
		{
			name:            "Request index.html in a nested directory indirectly",
			path:            "/static",
			requestPath:     "/static/nested/",
			expectedStatus:  http.StatusOK,
			expectedContent: "hello",
		},
	}

	subFS, err := fs.Sub(testdataFS, "testdata")
	assert.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := Middleware(tt.path, http.FS(subFS))
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAlreadyReported)
			})

			handlerToTest := middleware(nextHandler)

			for i := 0; i < 2; i++ {
				req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
				rr := httptest.NewRecorder()

				handlerToTest.ServeHTTP(rr, req)

				assert.Equal(t, tt.expectedStatus, rr.Code)
				if len(tt.expectedContent) > 0 {
					assert.Equal(t, tt.expectedContent, rr.Body.String())
				}
			}
		})
	}
}

func TestEnsureTrailingSlash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"path", "path/"},
		{"path/", "path/"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ensureTrailingSlash(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEnsureNoTrailingSlash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"path", "path"},
		{"path/", "path"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ensureNoTrailingSlash(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
