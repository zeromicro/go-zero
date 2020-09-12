package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/rest/internal"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestLogHandler(t *testing.T) {
	handlers := []func(handler http.Handler) http.Handler{
		LogHandler,
		DetailedLogHandler,
	}

	for _, logHandler := range handlers {
		req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
		handler := logHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Context().Value(internal.LogContext).(*internal.LogCollector).Append("anything")
			w.Header().Set("X-Test", "test")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, err := w.Write([]byte("content"))
			assert.Nil(t, err)
		}))

		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
		assert.Equal(t, "test", resp.Header().Get("X-Test"))
		assert.Equal(t, "content", resp.Body.String())
	}
}

func TestLogHandlerSlow(t *testing.T) {
	handlers := []func(handler http.Handler) http.Handler{
		LogHandler,
		DetailedLogHandler,
	}

	for _, logHandler := range handlers {
		req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
		handler := logHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(slowThreshold + time.Millisecond*50)
		}))

		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	}
}

func BenchmarkLogHandler(b *testing.B) {
	b.ReportAllocs()

	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	handler := LogHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < b.N; i++ {
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
	}
}
