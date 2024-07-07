package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zeromicro/go-zero/core/metainfo"
)

func TestCustomKeysHandler(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	testKey := metainfo.PrefixPass + "test"
	testCustomKeysData := map[string]string{testKey: testKey}

	handler := CustomKeysHandler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		md := metainfo.GetMapFromContext(request.Context())
		assert.Equal(t, testCustomKeysData, md)
		// todo(cong): test logger fields.
		wg.Done()
	}))

	req := httptest.NewRequest(http.MethodPost, "http://localhost",
		bytes.NewBufferString("123456789012345"))

	req.Header.Set(testKey, testKey)

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	wg.Wait()
}
