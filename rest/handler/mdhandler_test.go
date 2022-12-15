package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/md"
)

func TestMdHandler(t *testing.T) {
	handler := MdHandler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, []string{"v1", "v2"}, md.FromContext(request.Context()).Values("colors"))
	}))

	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.Header.Add("Colors", "v1")
	req.Header.Add("Colors", "v2")

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
}
