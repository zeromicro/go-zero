package response

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderOnceResponseWriter_Flush(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cw := NewHeaderOnceResponseWriter(w)
		cw.Header().Set("X-Test", "test")
		cw.WriteHeader(http.StatusServiceUnavailable)
		cw.WriteHeader(http.StatusExpectationFailed)
		_, err := cw.Write([]byte("content"))
		assert.Nil(t, err)

		flusher, ok := cw.(http.Flusher)
		assert.True(t, ok)
		flusher.Flush()
	})

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
	assert.Equal(t, "test", resp.Header().Get("X-Test"))
	assert.Equal(t, "content", resp.Body.String())
}

func TestHeaderOnceResponseWriter_Hijack(t *testing.T) {
	resp := httptest.NewRecorder()
	writer := &HeaderOnceResponseWriter{
		w: resp,
	}
	assert.NotPanics(t, func() {
		writer.Hijack()
	})

	writer = &HeaderOnceResponseWriter{
		w: mockedHijackable{resp},
	}
	assert.NotPanics(t, func() {
		writer.Hijack()
	})
}

type mockedHijackable struct {
	*httptest.ResponseRecorder
}

func (m mockedHijackable) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}
