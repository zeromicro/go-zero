package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/lang"
)

const conns = 4

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestMaxConnsHandler(t *testing.T) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(conns)
	done := make(chan lang.PlaceholderType)
	defer close(done)

	maxConns := MaxConns(conns)
	handler := maxConns(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		waitGroup.Done()
		<-done
	}))

	for i := 0; i < conns; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
			handler.ServeHTTP(httptest.NewRecorder(), req)
		}()
	}

	waitGroup.Wait()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
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

	maxConns := MaxConns(0)
	handler := maxConns(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get(key)
		if val == value {
			waitGroup.Done()
			<-done
		}
	}))

	for i := 0; i < conns; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
			req.Header.Set(key, value)
			handler.ServeHTTP(httptest.NewRecorder(), req)
		}()
	}

	waitGroup.Wait()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
