package mcp

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/rest/pathvar"
)

// RequestMetadata carries selected request-scoped values into MCP handlers.
type RequestMetadata struct {
	Headers map[string][]string
	Query   map[string][]string
	Path    map[string]string
}

type requestMetadataCtxKey struct{}

// RequestMetadataFromContext returns metadata extracted at the transport boundary.
func RequestMetadataFromContext(ctx context.Context) (RequestMetadata, bool) {
	metadata, ok := ctx.Value(requestMetadataCtxKey{}).(RequestMetadata)
	if !ok {
		return RequestMetadata{}, false
	}

	return normalizeRequestMetadata(metadata), true
}

// HeaderFromContext returns the first header value for key.
func HeaderFromContext(ctx context.Context, key string) (string, bool) {
	metadata, ok := RequestMetadataFromContext(ctx)
	if !ok {
		return "", false
	}

	vals := metadata.Headers[http.CanonicalHeaderKey(key)]
	if len(vals) == 0 {
		return "", false
	}

	return vals[0], true
}

// QueryFromContext returns the first query value for key.
func QueryFromContext(ctx context.Context, key string) (string, bool) {
	metadata, ok := RequestMetadataFromContext(ctx)
	if !ok {
		return "", false
	}

	vals := metadata.Query[key]
	if len(vals) == 0 {
		return "", false
	}

	return vals[0], true
}

// PathFromContext returns the path variable value for key.
func PathFromContext(ctx context.Context, key string) (string, bool) {
	metadata, ok := RequestMetadataFromContext(ctx)
	if !ok {
		return "", false
	}

	val, ok := metadata.Path[key]
	if !ok {
		return "", false
	}

	return val, true
}

// DefaultRequestMetadataExtractor extracts headers, query values, and path variables.
func DefaultRequestMetadataExtractor(r *http.Request) RequestMetadata {
	metadata := RequestMetadata{
		Headers: make(map[string][]string, len(r.Header)),
		Query:   make(map[string][]string),
		Path:    clonePathVars(pathvar.Vars(r)),
	}

	for key, vals := range r.Header {
		metadata.Headers[http.CanonicalHeaderKey(key)] = append([]string(nil), vals...)
	}

	if r.URL != nil {
		for key, vals := range r.URL.Query() {
			metadata.Query[key] = append([]string(nil), vals...)
		}
	}

	return metadata
}

func normalizeRequestMetadata(metadata RequestMetadata) RequestMetadata {
	return RequestMetadata{
		Headers: cloneHeaderValues(metadata.Headers),
		Query:   cloneHeaderValues(metadata.Query),
		Path:    clonePathVars(metadata.Path),
	}
}

func cloneHeaderValues(values map[string][]string) map[string][]string {
	if len(values) == 0 {
		return nil
	}

	cloned := make(map[string][]string, len(values))
	for key, vals := range values {
		cloned[key] = append([]string(nil), vals...)
	}

	return cloned
}

func clonePathVars(values map[string]string) map[string]string {
	if len(values) == 0 {
		return nil
	}

	cloned := make(map[string]string, len(values))
	for key, val := range values {
		cloned[key] = val
	}

	return cloned
}
