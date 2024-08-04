package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/metainfo"
)

func TestCustomKeysHandler(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	testKey := metainfo.PrefixPass + "test"
	testCustomKeysData := map[string]string{testKey: testKey}

	handler := CustomKeysHandler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		t.Helper() // Mark this as a test helper
		md := metainfo.GetMapFromContext(request.Context())
		assert.Equal(t, testCustomKeysData, md)
		writer.WriteHeader(http.StatusOK)
		wg.Done()
	}))

	req := httptest.NewRequest(http.MethodPost, "http://localhost",
		bytes.NewBufferString("123456789012345"))

	req.Header.Set(testKey, testKey)

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	// Use a context with timeout to avoid indefinite blocking
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Empty(t, resp.Body.String()) // Check if the response body is empty as expected
	case <-ctx.Done():
		t.Fatal("Test timed out")
	}
}
