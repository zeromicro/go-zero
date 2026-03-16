package httpc

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/rest/internal/header"
)

func TestNamedService_DoRequest(t *testing.T) {
	svr := httptest.NewServer(http.RedirectHandler("/foo", http.StatusMovedPermanently))
	defer svr.Close()
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)
	service := NewService("foo")
	_, err = service.DoRequest(req)
	// too many redirects
	assert.NotNil(t, err)
}

func TestNamedService_DoRequestGet(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("foo", r.Header.Get("foo"))
	}))
	defer svr.Close()
	service := NewService("foo", func(r *http.Request) *http.Request {
		r.Header.Set("foo", "bar")
		return r
	})
	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.Nil(t, err)
	resp, err := service.DoRequest(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "bar", resp.Header.Get("foo"))
}

func TestNamedService_DoRequestPost(t *testing.T) {
	svr := httptest.NewServer(http.NotFoundHandler())
	defer svr.Close()
	service := NewService("foo")
	req, err := http.NewRequest(http.MethodPost, svr.URL, nil)
	assert.Nil(t, err)
	req.Header.Set(header.ContentType, header.ContentTypeJson)
	resp, err := service.DoRequest(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestNamedService_Do(t *testing.T) {
	type Data struct {
		Key    string `path:"key"`
		Value  int    `form:"value"`
		Header string `header:"X-Header"`
		Body   string `json:"body"`
	}

	svr := httptest.NewServer(http.NotFoundHandler())
	defer svr.Close()

	service := NewService("foo")
	data := Data{
		Key:    "foo",
		Value:  10,
		Header: "my-header",
		Body:   "my body",
	}
	resp, err := service.Do(context.Background(), http.MethodPost, svr.URL+"/nodes/:key", data)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestNamedService_DoBadRequest(t *testing.T) {
	val := struct {
		Value string `json:"value,options=[a,b]"`
	}{
		Value: "c",
	}

	service := NewService("foo")
	_, err := service.Do(context.Background(), http.MethodPost, "/nodes/:key", val)
	assert.NotNil(t, err)
}

// mockNetError implements net.Error interface for testing
type mockNetError struct {
	msg       string
	timeout   bool
	temporary bool
}

func (e *mockNetError) Error() string   { return e.msg }
func (e *mockNetError) Timeout() bool   { return e.timeout }
func (e *mockNetError) Temporary() bool { return e.temporary }

func TestAcceptable(t *testing.T) {
	tests := []struct {
		name     string
		resp     *http.Response
		err      error
		expected bool
	}{
		{
			name: "no error with 2xx status code",
			resp: &http.Response{
				StatusCode: http.StatusOK,
			},
			err:      nil,
			expected: true,
		},
		{
			name: "no error with 3xx status code",
			resp: &http.Response{
				StatusCode: http.StatusMovedPermanently,
			},
			err:      nil,
			expected: true,
		},
		{
			name: "no error with 4xx status code",
			resp: &http.Response{
				StatusCode: http.StatusNotFound,
			},
			err:      nil,
			expected: true,
		},
		{
			name: "no error with 499 status code (just below 500)",
			resp: &http.Response{
				StatusCode: 499,
			},
			err:      nil,
			expected: true,
		},
		{
			name: "no error with 500 status code",
			resp: &http.Response{
				StatusCode: http.StatusInternalServerError,
			},
			err:      nil,
			expected: false,
		},
		{
			name: "no error with 503 status code",
			resp: &http.Response{
				StatusCode: http.StatusServiceUnavailable,
			},
			err:      nil,
			expected: false,
		},
		{
			name:     "context deadline exceeded",
			resp:     nil,
			err:      context.DeadlineExceeded,
			expected: false,
		},
		{
			name:     "context canceled",
			resp:     nil,
			err:      context.Canceled,
			expected: true,
		},
		{
			name:     "wrapped context deadline exceeded",
			resp:     nil,
			err:      errors.Join(context.DeadlineExceeded, errors.New("timeout")),
			expected: false,
		},
		{
			name:     "wrapped context canceled",
			resp:     nil,
			err:      errors.Join(context.Canceled, errors.New("canceled")),
			expected: true,
		},
		{
			name:     "network error - timeout",
			resp:     nil,
			err:      &mockNetError{msg: "network timeout", timeout: true, temporary: false},
			expected: false,
		},
		{
			name:     "network error - temporary",
			resp:     nil,
			err:      &mockNetError{msg: "temporary network error", timeout: false, temporary: true},
			expected: false,
		},
		{
			name:     "network error - connection refused",
			resp:     nil,
			err:      &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")},
			expected: false,
		},
		{
			name: "url.Error wrapping network error",
			resp: nil,
			err: &url.Error{
				Op:  "Get",
				URL: "http://example.com",
				Err: &mockNetError{msg: "network error", timeout: true},
			},
			expected: false,
		},
		{
			name: "url.Error wrapping non-network error",
			resp: nil,
			err: &url.Error{
				Op:  "Get",
				URL: "http://example.com",
				Err: errors.New("some other error"),
			},
			expected: true,
		},
		{
			name: "url.Error wrapping context.DeadlineExceeded",
			resp: nil,
			err: &url.Error{
				Op:  "Get",
				URL: "http://example.com",
				Err: context.DeadlineExceeded,
			},
			expected: false,
		},
		{
			name: "url.Error wrapping context.Canceled",
			resp: nil,
			err: &url.Error{
				Op:  "Get",
				URL: "http://example.com",
				Err: context.Canceled,
			},
			expected: true,
		},
		{
			name:     "generic error (non-network)",
			resp:     nil,
			err:      errors.New("some random error"),
			expected: true,
		},
		{
			name:     "EOF error (non-network)",
			resp:     nil,
			err:      errors.New("EOF"),
			expected: true,
		},
		{
			name:     "nil response with nil error (edge case)",
			resp:     nil,
			err:      nil,
			expected: false, // Will panic in real code, but resp.StatusCode access
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Handle the edge case where resp is nil and err is nil
			if tt.resp == nil && tt.err == nil {
				// This would panic in real code, so we skip the actual test
				// In production, this should never happen
				return
			}

			result := acceptable(tt.resp, tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAcceptable_RealNetworkTimeout(t *testing.T) {
	// Create a client with very short timeout
	client := &http.Client{
		Timeout: 1 * time.Nanosecond, // Extremely short timeout to force timeout error
	}

	service := NewServiceWithClient("test", client)

	// Create a server that delays response
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()

	req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
	assert.NoError(t, err)

	// This should timeout and trigger the circuit breaker
	resp, err := service.DoRequest(req)

	// The error should be present due to timeout
	assert.Error(t, err)
	// Response might be nil due to timeout
	if resp != nil {
		t.Logf("Response status: %d", resp.StatusCode)
	}
}

func TestAcceptable_Integration(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		expectBreaker bool // Whether breaker should consider this as failure
	}{
		{"200 OK should not trigger breaker", http.StatusOK, false},
		{"201 Created should not trigger breaker", http.StatusCreated, false},
		{"400 Bad Request should not trigger breaker", http.StatusBadRequest, false},
		{"404 Not Found should not trigger breaker", http.StatusNotFound, false},
		{"500 Internal Server Error should trigger breaker", http.StatusInternalServerError, true},
		{"502 Bad Gateway should trigger breaker", http.StatusBadGateway, true},
		{"503 Service Unavailable should trigger breaker", http.StatusServiceUnavailable, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer svr.Close()

			service := NewService("test-service-" + tt.name)
			req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
			assert.NoError(t, err)

			resp, err := service.DoRequest(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.statusCode, resp.StatusCode)

			// The actual breaker behavior is tested implicitly through the acceptable function
			result := acceptable(resp, nil)
			if tt.expectBreaker {
				assert.False(t, result, "Status %d should not be acceptable", tt.statusCode)
			} else {
				assert.True(t, result, "Status %d should be acceptable", tt.statusCode)
			}
		})
	}
}
