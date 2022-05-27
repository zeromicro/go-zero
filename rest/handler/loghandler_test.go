package handler

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/rest/internal"
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

			flusher, ok := w.(http.Flusher)
			assert.True(t, ok)
			flusher.Flush()
		}))

		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
		assert.Equal(t, "test", resp.Header().Get("X-Test"))
		assert.Equal(t, "content", resp.Body.String())
	}
}

func TestLogHandlerVeryLong(t *testing.T) {
	var buf bytes.Buffer
	for i := 0; i < limitBodyBytes<<1; i++ {
		buf.WriteByte('a')
	}

	req := httptest.NewRequest(http.MethodPost, "http://localhost", &buf)
	handler := LogHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Context().Value(internal.LogContext).(*internal.LogCollector).Append("anything")
		io.Copy(ioutil.Discard, r.Body)
		w.Header().Set("X-Test", "test")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, err := w.Write([]byte("content"))
		assert.Nil(t, err)

		flusher, ok := w.(http.Flusher)
		assert.True(t, ok)
		flusher.Flush()
	}))

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
	assert.Equal(t, "test", resp.Header().Get("X-Test"))
	assert.Equal(t, "content", resp.Body.String())
}

func TestLogHandlerSlow(t *testing.T) {
	handlers := []func(handler http.Handler) http.Handler{
		LogHandler,
		DetailedLogHandler,
	}

	for _, logHandler := range handlers {
		req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
		handler := logHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(defaultSlowThreshold + time.Millisecond*50)
		}))

		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	}
}

func TestLogHandler_Hijack(t *testing.T) {
	resp := httptest.NewRecorder()
	writer := &loggedResponseWriter{
		w: resp,
	}
	assert.NotPanics(t, func() {
		writer.Hijack()
	})

	writer = &loggedResponseWriter{
		w: mockedHijackable{resp},
	}
	assert.NotPanics(t, func() {
		writer.Hijack()
	})
}

func TestDetailedLogHandler_Hijack(t *testing.T) {
	resp := httptest.NewRecorder()
	writer := &detailLoggedResponseWriter{
		writer: &loggedResponseWriter{
			w: resp,
		},
	}
	assert.NotPanics(t, func() {
		writer.Hijack()
	})

	writer = &detailLoggedResponseWriter{
		writer: &loggedResponseWriter{
			w: mockedHijackable{resp},
		},
	}
	assert.NotPanics(t, func() {
		writer.Hijack()
	})
}

func TestSetSlowThreshold(t *testing.T) {
	assert.Equal(t, defaultSlowThreshold, slowThreshold.Load())
	SetSlowThreshold(time.Second)
	assert.Equal(t, time.Second, slowThreshold.Load())
}

func TestWrapMethodWithColor(t *testing.T) {
	// no tty
	assert.Equal(t, http.MethodGet, wrapMethod(http.MethodGet))
	assert.Equal(t, http.MethodPost, wrapMethod(http.MethodPost))
	assert.Equal(t, http.MethodPut, wrapMethod(http.MethodPut))
	assert.Equal(t, http.MethodDelete, wrapMethod(http.MethodDelete))
	assert.Equal(t, http.MethodPatch, wrapMethod(http.MethodPatch))
	assert.Equal(t, http.MethodHead, wrapMethod(http.MethodHead))
	assert.Equal(t, http.MethodOptions, wrapMethod(http.MethodOptions))
	assert.Equal(t, http.MethodConnect, wrapMethod(http.MethodConnect))
	assert.Equal(t, http.MethodTrace, wrapMethod(http.MethodTrace))
}

func TestWrapStatusCodeWithColor(t *testing.T) {
	// no tty
	assert.Equal(t, "200", wrapStatusCode(http.StatusOK))
	assert.Equal(t, "302", wrapStatusCode(http.StatusFound))
	assert.Equal(t, "404", wrapStatusCode(http.StatusNotFound))
	assert.Equal(t, "500", wrapStatusCode(http.StatusInternalServerError))
	assert.Equal(t, "503", wrapStatusCode(http.StatusServiceUnavailable))
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
