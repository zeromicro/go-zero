package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxBytesHandler(t *testing.T) {
	maxb := MaxBytesHandler(10)
	handler := maxb(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest(http.MethodPost, "http://localhost",
		bytes.NewBufferString("123456789012345"))
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusRequestEntityTooLarge, resp.Code)

	req = httptest.NewRequest(http.MethodPost, "http://localhost", bytes.NewBufferString("12345"))
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestMaxBytesHandlerNoLimit(t *testing.T) {
	maxb := MaxBytesHandler(-1)
	handler := maxb(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest(http.MethodPost, "http://localhost",
		bytes.NewBufferString("123456789012345"))
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
