package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/trace/tracespec"
)

func TestTracingHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.Header.Set("X-Trace-ID", "theid")
	handler := TracingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, ok := r.Context().Value(tracespec.TracingKey).(tracespec.Trace)
		assert.True(t, ok)
		assert.Equal(t, "theid", span.TraceId())
	}))

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
