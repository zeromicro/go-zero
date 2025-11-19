package handler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx/logtest"
	"github.com/zeromicro/go-zero/rest/internal"
	"github.com/zeromicro/go-zero/rest/internal/response"
)

func TestLogHandler(t *testing.T) {
	handlers := []func(handler http.Handler) http.Handler{
		LogHandler,
		DetailedLogHandler,
	}

	for _, logHandler := range handlers {
		req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
		handler := logHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			internal.LogCollectorFromContext(r.Context()).Append("anything")
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
		internal.LogCollectorFromContext(r.Context()).Append("anything")
		_, _ = io.Copy(io.Discard, r.Body)
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
		req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
		handler := logHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(defaultSlowThreshold + time.Millisecond*50)
		}))

		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	}
}

func TestLogHandlerSSE(t *testing.T) {
	handlers := []func(handler http.Handler) http.Handler{
		LogHandler,
		DetailedLogHandler,
	}

	for _, logHandler := range handlers {
		t.Run("SSE request with normal duration", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
			req.Header.Set(headerAccept, valueSSE)

			handler := logHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(defaultSlowThreshold + time.Second)
				w.WriteHeader(http.StatusOK)
			}))

			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusOK, resp.Code)
		})

		t.Run("SSE request exceeding SSE threshold", func(t *testing.T) {
			originalThreshold := sseSlowThreshold.Load()
			SetSSESlowThreshold(time.Millisecond * 100)
			defer SetSSESlowThreshold(originalThreshold)

			req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
			req.Header.Set(headerAccept, valueSSE)

			handler := logHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(time.Millisecond * 150)
				w.WriteHeader(http.StatusOK)
			}))

			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
	}
}

func TestLogHandlerThresholdSelection(t *testing.T) {
	tests := []struct {
		name          string
		acceptHeader  string
		expectedIsSSE bool
	}{
		{
			name:          "Regular HTTP request",
			acceptHeader:  "text/html",
			expectedIsSSE: false,
		},
		{
			name:          "SSE request",
			acceptHeader:  valueSSE,
			expectedIsSSE: true,
		},
		{
			name:          "No Accept header",
			acceptHeader:  "",
			expectedIsSSE: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
			if tt.acceptHeader != "" {
				req.Header.Set(headerAccept, tt.acceptHeader)
			}

			SetSlowThreshold(time.Millisecond * 100)
			SetSSESlowThreshold(time.Millisecond * 200)
			defer func() {
				SetSlowThreshold(defaultSlowThreshold)
				SetSSESlowThreshold(defaultSSESlowThreshold)
			}()

			handler := LogHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(time.Millisecond * 150)
				w.WriteHeader(http.StatusOK)
			}))

			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusOK, resp.Code)
		})
	}
}

func TestDetailedLogHandler_LargeBody(t *testing.T) {
	lbuf := logtest.NewCollector(t)

	var buf bytes.Buffer
	for i := 0; i < limitDetailedBodyBytes<<2; i++ {
		buf.WriteByte('a')
	}

	req := httptest.NewRequest(http.MethodPost, "http://localhost", &buf)
	h := DetailedLogHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
	}))
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)

	// extra 200 for the length of POST request headers
	assert.True(t, len(lbuf.Content()) < limitDetailedBodyBytes+200)
}

func TestDetailedLogHandler_Hijack(t *testing.T) {
	resp := httptest.NewRecorder()
	writer := &detailLoggedResponseWriter{
		writer: response.NewWithCodeResponseWriter(resp),
	}
	assert.NotPanics(t, func() {
		_, _, _ = writer.Hijack()
	})

	writer = &detailLoggedResponseWriter{
		writer: response.NewWithCodeResponseWriter(resp),
	}
	assert.NotPanics(t, func() {
		_, _, _ = writer.Hijack()
	})

	writer = &detailLoggedResponseWriter{
		writer: response.NewWithCodeResponseWriter(mockedHijackable{
			ResponseRecorder: resp,
		}),
	}
	assert.NotPanics(t, func() {
		_, _, _ = writer.Hijack()
	})
}

func TestSetSlowThreshold(t *testing.T) {
	assert.Equal(t, defaultSlowThreshold, slowThreshold.Load())
	SetSlowThreshold(time.Second)
	assert.Equal(t, time.Second, slowThreshold.Load())
}

func TestSetSSESlowThreshold(t *testing.T) {
	assert.Equal(t, defaultSSESlowThreshold, sseSlowThreshold.Load())
	SetSSESlowThreshold(time.Minute * 10)
	assert.Equal(t, time.Minute*10, sseSlowThreshold.Load())
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

func TestDumpRequest(t *testing.T) {
	const errMsg = "error"
	r := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
	r.Body = mockedReadCloser{errMsg: errMsg}
	assert.Equal(t, errMsg, dumpRequest(r))
}

func BenchmarkLogHandler(b *testing.B) {
	b.ReportAllocs()

	req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
	handler := LogHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < b.N; i++ {
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
	}
}

type mockedReadCloser struct {
	errMsg string
}

func (m mockedReadCloser) Read(_ []byte) (n int, err error) {
	return 0, errors.New(m.errMsg)
}

func (m mockedReadCloser) Close() error {
	return nil
}
