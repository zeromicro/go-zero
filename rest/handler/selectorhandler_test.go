package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/selector"
)

func TestSelectorHandler(t *testing.T) {
	handler := SelectorHandler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, []string{"v1", "v2"}, selector.ColorsFromContext(request.Context()).Colors())
		assert.Equal(t, selector.DefaultSelector, selector.SelectorFromContext(request.Context()))
	}))

	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.Header.Add("Colors", "v1")
	req.Header.Add("Colors", "v2")
	req.Header.Add("Selector", selector.DefaultSelector)

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
}
