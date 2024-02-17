package internal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx/logtest"
)

func TestInfo(t *testing.T) {
	collector := new(LogCollector)
	req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
	req = req.WithContext(WithLogCollector(req.Context(), collector))
	Info(req, "first")
	Infof(req, "second %s", "third")
	val := collector.Flush()
	assert.True(t, strings.Contains(val, "first"))
	assert.True(t, strings.Contains(val, "second"))
	assert.True(t, strings.Contains(val, "third"))
	assert.True(t, strings.Contains(val, "\n"))
}

func TestError(t *testing.T) {
	c := logtest.NewCollector(t)
	req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
	Error(req, "first")
	Errorf(req, "second %s", "third")
	val := c.String()
	assert.True(t, strings.Contains(val, "first"))
	assert.True(t, strings.Contains(val, "second"))
	assert.True(t, strings.Contains(val, "third"))
}

func TestLogCollectorContext(t *testing.T) {
	ctx := context.Background()
	assert.Nil(t, LogCollectorFromContext(ctx))
	collector := new(LogCollector)
	ctx = WithLogCollector(ctx, collector)
	assert.Equal(t, collector, LogCollectorFromContext(ctx))
}
