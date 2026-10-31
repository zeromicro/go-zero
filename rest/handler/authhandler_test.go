package handler

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
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

type countingReadCloser struct {
	io.Reader
	total int64
}

func (c *countingReadCloser) Read(p []byte) (int, error) {
	n, err := c.Reader.Read(p)
	c.total += int64(n)
	return n, err
}

func (c *countingReadCloser) Close() error {
	return nil
}

// TestAuthHandlerFailedWithGzipBody guards against GHSA-h4vj-585c-2r8x: a
// failed JWT authentication must not force the server to read/decompress
// the full (potentially huge) gzip-encoded request body just to log it.
func TestAuthHandlerFailedWithGzipBody(t *testing.T) {
	// highly compressible payload, so a small compressed input unpacks to
	// a payload far larger than the detailed-log limit.
	plain := bytes.Repeat([]byte("A"), 16<<20) // 16 MiB decompressed
	var packed bytes.Buffer
	zw := gzip.NewWriter(&packed)
	_, err := zw.Write(plain)
	assert.Nil(t, err)
	assert.Nil(t, zw.Close())

	counter := &countingReadCloser{Reader: bytes.NewReader(packed.Bytes())}
	req := httptest.NewRequest(http.MethodPost, "http://localhost", nil)
	req.Body = counter
	req.Header.Set("Content-Encoding", "gzip")

	handler := GunzipHandler(Authorize("B63F477D-BBA3-4E52-96D3-C0034C27694A")(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})))

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	// the compressed bytes actually consumed must stay far below the full
	// compressed payload size, proving the body wasn't fully decompressed.
	assert.Less(t, counter.total, int64(packed.Len()))
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
