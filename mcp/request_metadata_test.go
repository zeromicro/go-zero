package mcp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func TestDefaultRequestMetadataExtractor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/sse?tenant=t1&trace=abc", nil)
	req.Header.Add("X-Tenant-Id", "tenant-from-header")
	req = pathvar.WithVars(req, map[string]string{"tool": "sum"})

	metadata := DefaultRequestMetadataExtractor(req)
	header, ok := metadata.Headers["X-Tenant-Id"]
	assert.True(t, ok)
	assert.Equal(t, []string{"tenant-from-header"}, header)
	assert.Equal(t, []string{"t1"}, metadata.Query["tenant"])
	assert.Equal(t, "sum", metadata.Path["tool"])
}

func TestRequestMetadataContextHelpers(t *testing.T) {
	ctx := context.WithValue(context.Background(), requestMetadataCtxKey{}, RequestMetadata{
		Headers: map[string][]string{"X-Trace-Id": {"trace-1"}},
		Query:   map[string][]string{"tenant": {"foo"}},
		Path:    map[string]string{"scope": "prod"},
	})

	metadata, ok := RequestMetadataFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, []string{"trace-1"}, metadata.Headers["X-Trace-Id"])

	header, ok := HeaderFromContext(ctx, "x-trace-id")
	assert.True(t, ok)
	assert.Equal(t, "trace-1", header)

	query, ok := QueryFromContext(ctx, "tenant")
	assert.True(t, ok)
	assert.Equal(t, "foo", query)

	path, ok := PathFromContext(ctx, "scope")
	assert.True(t, ok)
	assert.Equal(t, "prod", path)
}

func TestRequestMetadataFromContextNotFound(t *testing.T) {
	_, ok := RequestMetadataFromContext(context.Background())
	assert.False(t, ok)

	_, ok = HeaderFromContext(context.Background(), "x-test")
	assert.False(t, ok)

	_, ok = QueryFromContext(context.Background(), "tenant")
	assert.False(t, ok)

	_, ok = PathFromContext(context.Background(), "tenant")
	assert.False(t, ok)
}

func TestWrapRequestMetadata(t *testing.T) {
	s := &mcpServerImpl{
		options: serverOptions{
			requestMetadataExtractor: DefaultRequestMetadataExtractor,
		},
	}

	called := false
	handler := s.wrapRequestMetadata(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		called = true
		header, ok := HeaderFromContext(r.Context(), "x-tenant-id")
		assert.True(t, ok)
		assert.Equal(t, "tenant-1", header)

		query, ok := QueryFromContext(r.Context(), "tenant")
		assert.True(t, ok)
		assert.Equal(t, "q-tenant", query)
	}))

	req := httptest.NewRequest(http.MethodGet, "/sse?tenant=q-tenant", nil)
	req.Header.Set("X-Tenant-Id", "tenant-1")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.True(t, called)
}

func TestWrapRequestMetadataNoExtractor(t *testing.T) {
	s := &mcpServerImpl{}

	called := false
	handler := s.wrapRequestMetadata(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		called = true
		_, ok := RequestMetadataFromContext(r.Context())
		assert.False(t, ok)
	}))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/sse", nil))

	assert.True(t, called)
}

func TestWrapRequestMetadataCanonicalizesCustomHeaders(t *testing.T) {
	s := &mcpServerImpl{
		options: serverOptions{
			requestMetadataExtractor: func(*http.Request) RequestMetadata {
				return RequestMetadata{
					Headers: map[string][]string{
						"x-tenant-id": {"tenant-lower"},
					},
				}
			},
		},
	}

	called := false
	handler := s.wrapRequestMetadata(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		called = true
		header, ok := HeaderFromContext(r.Context(), "X-Tenant-Id")
		assert.True(t, ok)
		assert.Equal(t, "tenant-lower", header)
	}))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/sse", nil))

	assert.True(t, called)
}

func TestRequestMetadataFromContextReturnsCopy(t *testing.T) {
	ctx := context.WithValue(context.Background(), requestMetadataCtxKey{}, RequestMetadata{
		Headers: map[string][]string{"X-Trace-Id": {"trace-1"}},
	})

	metadata, ok := RequestMetadataFromContext(ctx)
	assert.True(t, ok)
	metadata.Headers["X-Trace-Id"][0] = "mutated"
	metadata.Headers["X-New"] = []string{"new"}

	fresh, ok := RequestMetadataFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, []string{"trace-1"}, fresh.Headers["X-Trace-Id"])
	assert.Nil(t, fresh.Headers["X-New"])
}
