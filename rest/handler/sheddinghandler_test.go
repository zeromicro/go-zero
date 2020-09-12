package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/load"
	"github.com/tal-tech/go-zero/core/stat"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestSheddingHandlerAccept(t *testing.T) {
	metrics := stat.NewMetrics("unit-test")
	shedder := mockShedder{
		allow: true,
	}
	sheddingHandler := SheddingHandler(shedder, metrics)
	handler := sheddingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test", "test")
		_, err := w.Write([]byte("content"))
		assert.Nil(t, err)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.Header.Set("X-Test", "test")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "test", resp.Header().Get("X-Test"))
	assert.Equal(t, "content", resp.Body.String())
}

func TestSheddingHandlerFail(t *testing.T) {
	metrics := stat.NewMetrics("unit-test")
	shedder := mockShedder{
		allow: true,
	}
	sheddingHandler := SheddingHandler(shedder, metrics)
	handler := sheddingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
}

func TestSheddingHandlerReject(t *testing.T) {
	metrics := stat.NewMetrics("unit-test")
	shedder := mockShedder{
		allow: false,
	}
	sheddingHandler := SheddingHandler(shedder, metrics)
	handler := sheddingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
}

func TestSheddingHandlerNoShedding(t *testing.T) {
	metrics := stat.NewMetrics("unit-test")
	sheddingHandler := SheddingHandler(nil, metrics)
	handler := sheddingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

type mockShedder struct {
	allow bool
}

func (s mockShedder) Allow() (load.Promise, error) {
	if s.allow {
		return mockPromise{}, nil
	} else {
		return nil, load.ErrServiceOverloaded
	}
}

type mockPromise struct {
}

func (p mockPromise) Pass() {
}

func (p mockPromise) Fail() {
}
