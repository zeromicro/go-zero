package internal

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	collector := new(LogCollector)
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req = req.WithContext(context.WithValue(req.Context(), LogContext, collector))
	Info(req, "first")
	Infof(req, "second %s", "third")
	val := collector.Flush()
	assert.True(t, strings.Contains(val, "first"))
	assert.True(t, strings.Contains(val, "second"))
	assert.True(t, strings.Contains(val, "third"))
	assert.True(t, strings.Contains(val, "\n"))
}

func TestError(t *testing.T) {
	var writer strings.Builder
	log.SetOutput(&writer)
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	Error(req, "first")
	Errorf(req, "second %s", "third")
	val := writer.String()
	assert.True(t, strings.Contains(val, "first"))
	assert.True(t, strings.Contains(val, "second"))
	assert.True(t, strings.Contains(val, "third"))
	assert.True(t, strings.Contains(val, "\n"))
}
