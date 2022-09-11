package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/codec"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TestGunzipHandler(t *testing.T) {
	const message = "hello world"
	var wg sync.WaitGroup
	wg.Add(1)
	handler := GunzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.Equal(t, string(body), message)
		wg.Done()
	}))

	req := httptest.NewRequest(http.MethodPost, "http://localhost",
		bytes.NewReader(codec.Gzip([]byte(message))))
	req.Header.Set(httpx.ContentEncoding, gzipEncoding)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	wg.Wait()
}

func TestGunzipHandler_NoGzip(t *testing.T) {
	const message = "hello world"
	var wg sync.WaitGroup
	wg.Add(1)
	handler := GunzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.Equal(t, string(body), message)
		wg.Done()
	}))

	req := httptest.NewRequest(http.MethodPost, "http://localhost",
		strings.NewReader(message))
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	wg.Wait()
}

func TestGunzipHandler_NoGzipButTelling(t *testing.T) {
	const message = "hello world"
	handler := GunzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	req := httptest.NewRequest(http.MethodPost, "http://localhost",
		strings.NewReader(message))
	req.Header.Set(httpx.ContentEncoding, gzipEncoding)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
