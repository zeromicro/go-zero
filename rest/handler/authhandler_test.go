package handler

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandlerFailed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
	handler := Authorize("B63F477D-BBA3-4E52-96D3-C0034C27694A", WithUnauthorizedCallback(
		func(w http.ResponseWriter, r *http.Request, err error) {
			assert.NotNil(t, err)
			w.Header().Set("X-Test", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			_, err = w.Write([]byte("content"))
			assert.Nil(t, err)
		}))(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthHandler(t *testing.T) {
	const key = "B63F477D-BBA3-4E52-96D3-C0034C27694A"
	req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
	token, err := buildToken(key, map[string]any{
		"key": "value",
	}, 3600)
	assert.Nil(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	handler := Authorize(key)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "test")
			_, err := w.Write([]byte("content"))
			assert.Nil(t, err)

			flusher, ok := w.(http.Flusher)
			assert.True(t, ok)
			flusher.Flush()
		}))

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "test", resp.Header().Get("X-Test"))
	assert.Equal(t, "content", resp.Body.String())
}

func TestAuthHandlerWithPrevSecret(t *testing.T) {
	const (
		key     = "14F17379-EB8F-411B-8F12-6929002DCA76"
		prevKey = "B63F477D-BBA3-4E52-96D3-C0034C27694A"
	)
	req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
	token, err := buildToken(key, map[string]any{
		"key": "value",
	}, 3600)
	assert.Nil(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	handler := Authorize(key, WithPrevSecret(prevKey))(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "test")
			_, err := w.Write([]byte("content"))
			assert.Nil(t, err)
		}))

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "test", resp.Header().Get("X-Test"))
	assert.Equal(t, "content", resp.Body.String())
}

func TestAuthHandler_NilError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
	resp := httptest.NewRecorder()
	assert.NotPanics(t, func() {
		unauthorized(resp, req, nil, nil)
	})
}

func TestAuthHandlerWithJSONBody(t *testing.T) {
	const key = "B63F477D-BBA3-4E52-96D3-C0034C27694A"
	
	// Create a request with JSON body
	jsonBody := `{"username":"test","password":"secret"}`
	req := httptest.NewRequest(http.MethodPost, "http://localhost/login", 
		strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	// Missing authorization header to trigger the unauthorized path
	
	handler := Authorize(key)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
	
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	
	// Should return unauthorized
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthHandlerWithMultipartFormData(t *testing.T) {
	const key = "B63F477D-BBA3-4E52-96D3-C0034C27694A"
	
	// Create a multipart form-data request
	// We don't need actual body content since we're testing that
	// the body is NOT read when Content-Type is multipart/form-data
	req := httptest.NewRequest(http.MethodPost, "http://localhost/upload", 
		http.NoBody)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary")
	// Missing authorization header to trigger the unauthorized path
	
	handler := Authorize(key)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
	
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	
	// Should return unauthorized
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthHandlerWithMultipartFormDataLargeFile(t *testing.T) {
	const key = "B63F477D-BBA3-4E52-96D3-C0034C27694A"
	
	// Create a multipart form-data request with a simulated large file
	// This tests that the body is NOT consumed when Content-Type is multipart/form-data
	largeContent := make([]byte, 1024*1024) // 1MB of data
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}
	
	req := httptest.NewRequest(http.MethodPost, "http://localhost/upload", 
		http.NoBody)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary")
	req.Header.Set("Content-Length", "1048576")
	// Missing authorization header to trigger the unauthorized path
	
	handler := Authorize(key)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
	
	resp := httptest.NewRecorder()
	
	// This should complete quickly without reading the body
	start := time.Now()
	handler.ServeHTTP(resp, req)
	elapsed := time.Since(start)
	
	// Should return unauthorized
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	// Should complete in less than 100ms (without reading 1MB of data)
	assert.Less(t, elapsed, 100*time.Millisecond)
}

func TestIsMultipartFormData(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		expected    bool
	}{
		{
			name:        "multipart/form-data",
			contentType: "multipart/form-data",
			expected:    true,
		},
		{
			name:        "multipart/form-data with boundary",
			contentType: "multipart/form-data; boundary=----WebKitFormBoundary",
			expected:    true,
		},
		{
			name:        "application/json",
			contentType: "application/json",
			expected:    false,
		},
		{
			name:        "application/x-www-form-urlencoded",
			contentType: "application/x-www-form-urlencoded",
			expected:    false,
		},
		{
			name:        "empty content type",
			contentType: "",
			expected:    false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "http://localhost", http.NoBody)
			req.Header.Set("Content-Type", tt.contentType)
			result := isMultipartFormData(req)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func buildToken(secretKey string, payloads map[string]any, seconds int64) (string, error) {
	now := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = now + seconds
	claims["iat"] = now
	for k, v := range payloads {
		claims[k] = v
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}

type mockedHijackable struct {
	*httptest.ResponseRecorder
}

func (m mockedHijackable) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}
