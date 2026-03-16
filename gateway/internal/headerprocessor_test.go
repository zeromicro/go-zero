package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildHeadersNoValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/", http.NoBody)
	req.Header.Add("a", "b")
	assert.Nil(t, ProcessHeaders(req.Header))
}

func TestBuildHeadersWithValues(t *testing.T) {
	req := httptest.NewRequest("GET", "/", http.NoBody)
	req.Header.Add("grpc-metadata-a", "b")
	req.Header.Add("grpc-metadata-b", "b")
	assert.ElementsMatch(t, []string{"gateway-a:b", "gateway-b:b"}, ProcessHeaders(req.Header))
}

func TestProcessHeadersWithTraceContext(t *testing.T) {
	req := httptest.NewRequest("GET", "/", http.NoBody)
	req.Header.Set("traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
	req.Header.Set("tracestate", "key1=value1,key2=value2")
	req.Header.Set("baggage", "userId=alice,serverNode=DF:28")

	headers := ProcessHeaders(req.Header)

	assert.Len(t, headers, 3)
	assert.Contains(t, headers, "traceparent:00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
	assert.Contains(t, headers, "tracestate:key1=value1,key2=value2")
	assert.Contains(t, headers, "baggage:userId=alice,serverNode=DF:28")
}

func TestProcessHeadersWithMixedHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/", http.NoBody)
	req.Header.Set("traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
	req.Header.Set("grpc-metadata-custom", "value1")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("tracestate", "key1=value1")

	headers := ProcessHeaders(req.Header)

	// Should include trace headers and grpc-metadata headers, but not regular headers
	assert.Len(t, headers, 3)
	assert.Contains(t, headers, "traceparent:00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
	assert.Contains(t, headers, "tracestate:key1=value1")
	assert.Contains(t, headers, "gateway-custom:value1")
}

func TestProcessHeadersTraceparentCaseInsensitive(t *testing.T) {
	tests := []struct {
		name        string
		headerKey   string
		headerVal   string
		expectedKey string
	}{
		{
			name:        "lowercase traceparent",
			headerKey:   "traceparent",
			headerVal:   "00-trace-span-01",
			expectedKey: "traceparent",
		},
		{
			name:        "uppercase Traceparent",
			headerKey:   "Traceparent",
			headerVal:   "00-trace-span-01",
			expectedKey: "traceparent",
		},
		{
			name:        "mixed case TraceParent",
			headerKey:   "TraceParent",
			headerVal:   "00-trace-span-01",
			expectedKey: "traceparent",
		},
		{
			name:        "lowercase tracestate",
			headerKey:   "tracestate",
			headerVal:   "key=value",
			expectedKey: "tracestate",
		},
		{
			name:        "mixed case TraceState",
			headerKey:   "TraceState",
			headerVal:   "key=value",
			expectedKey: "tracestate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", http.NoBody)
			req.Header.Set(tt.headerKey, tt.headerVal)

			headers := ProcessHeaders(req.Header)

			assert.Len(t, headers, 1)
			assert.Contains(t, headers, tt.expectedKey+":"+tt.headerVal)
		})
	}
}

func TestProcessHeadersEmptyHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/", http.NoBody)
	headers := ProcessHeaders(req.Header)
	assert.Empty(t, headers)
}
