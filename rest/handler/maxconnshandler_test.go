package handler

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/lang"
)

const conns = 4

func TestMaxConnsHandler(t *testing.T) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(conns)
	done := make(chan lang.PlaceholderType)
	defer close(done)

	maxConns := MaxConnsHandler(conns)
	handler := maxConns(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		waitGroup.Done()
		<-done
	}))

	for i := 0; i < conns; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
			handler.ServeHTTP(httptest.NewRecorder(), req)
		}()
	}

	waitGroup.Wait()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
}

func TestWithoutMaxConnsHandler(t *testing.T) {
	const (
		key   = "block"
		value = "1"
	)
	var waitGroup sync.WaitGroup
	waitGroup.Add(conns)
	done := make(chan lang.PlaceholderType)
	defer close(done)

	maxConns := MaxConnsHandler(0)
	handler := maxConns(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get(key)
		if val == value {
			waitGroup.Done()
			<-done
		}
	}))

	for i := 0; i < conns; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
			req.Header.Set(key, value)
			handler.ServeHTTP(httptest.NewRecorder(), req)
		}()
	}

	waitGroup.Wait()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
