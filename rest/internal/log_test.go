package internal

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/logx/logtest"
)

func TestInfo(t *testing.T) {
	collector := new(LogCollector)
	req := httptest.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
	req = req.WithContext(WithLogCollector(req.Context(), collector))
	Info(req, "id", "123456")
	Infof(req, "mobile_channel", "channel_%s", "first")
	val := collector.Flush()
	for _, field := range val {
		logx.Infow(fmt.Sprintf("%s_test", field.Key), field)
	}
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
