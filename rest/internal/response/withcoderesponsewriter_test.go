package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithCodeResponseWriter(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cw := &WithCodeResponseWriter{Writer: w}

		cw.Header().Set("X-Test", "test")
		cw.WriteHeader(http.StatusServiceUnavailable)
		assert.Equal(t, cw.Code, http.StatusServiceUnavailable)

		_, err := cw.Write([]byte("content"))
		assert.Nil(t, err)

		flusher, ok := http.ResponseWriter(cw).(http.Flusher)
		assert.True(t, ok)
		flusher.Flush()
	})

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
	assert.Equal(t, "test", resp.Header().Get("X-Test"))
	assert.Equal(t, "content", resp.Body.String())
}

func TestWithCodeResponseWriter_Hijack(t *testing.T) {
	resp := httptest.NewRecorder()
	writer := &WithCodeResponseWriter{
		Writer: resp,
	}
	assert.NotPanics(t, func() {
		writer.Hijack()
	})

	writer = &WithCodeResponseWriter{
		Writer: mockedHijackable{resp},
	}
	assert.NotPanics(t, func() {
		writer.Hijack()
	})
}
