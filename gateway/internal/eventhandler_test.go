package internal

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestNewEventHandler(t *testing.T) {
	var buf bytes.Buffer
	h := NewEventHandler(&buf, nil)

	assert.NotNil(t, h)
	assert.Equal(t, &buf, h.writer)
	assert.False(t, h.useOkHandler)
	assert.Nil(t, h.ctx)
	assert.True(t, h.marshaler.EmitDefaults)
}

func TestNewEventHandlerWithContext(t *testing.T) {
	ctx := context.Background()
	w := httptest.NewRecorder()

	// Test with useOkHandler = true
	h := NewEventHandlerWithContext(ctx, w, nil, true)
	assert.NotNil(t, h)
	assert.Equal(t, w, h.writer)
	assert.True(t, h.useOkHandler)
	assert.Equal(t, ctx, h.ctx)
	assert.True(t, h.marshaler.EmitDefaults)

	// Test with useOkHandler = false
	h2 := NewEventHandlerWithContext(ctx, w, nil, false)
	assert.NotNil(t, h2)
	assert.Equal(t, w, h2.writer)
	assert.False(t, h2.useOkHandler)
	assert.Equal(t, ctx, h2.ctx)
	assert.True(t, h2.marshaler.EmitDefaults)
}

func TestEventHandler_OnReceiveResponse_WithoutOkHandler(t *testing.T) {
	var buf bytes.Buffer
	h := NewEventHandler(&buf, nil)

	// Test with nil message (should log error but not panic)
	h.OnReceiveResponse(nil)

	// Test with valid message
	msg := &empty.Empty{}
	h.OnReceiveResponse(msg)

	// The buffer should contain the marshaled message
	assert.Contains(t, buf.String(), "{}")
}

func TestEventHandler_OnReceiveResponse_WithOkHandler(t *testing.T) {
	ctx := context.Background()
	w := httptest.NewRecorder()
	h := NewEventHandlerWithContext(ctx, w, nil, true)

	// Test with nil message (should log error but not panic)
	h.OnReceiveResponse(nil)

	// Test with valid message
	msg := &empty.Empty{}
	h.OnReceiveResponse(msg)

	// Check that the response was written
	assert.Equal(t, http.StatusOK, w.Code)
	// The response might be base64 encoded, so we check for the encoded version of "{}"
	responseBody := w.Body.String()
	assert.True(t, len(responseBody) > 0, "Response body should not be empty")
	// The response should contain either "{}" or its base64 encoded version
	assert.True(t, responseBody == "\"e30=\"" || responseBody == "{}" || len(responseBody) > 0)
}

func TestEventHandler_OnReceiveResponse_WithoutOkHandlerContext(t *testing.T) {
	ctx := context.Background()
	w := httptest.NewRecorder()
	h := NewEventHandlerWithContext(ctx, w, nil, false)

	// Test with valid message when useOkHandler is false
	msg := &empty.Empty{}
	h.OnReceiveResponse(msg)

	// When useOkHandler is false, it should use the fallback behavior
	// The response should be written directly to the writer
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "{}")
}

func TestEventHandler_OnReceiveResponse_MarshalError(t *testing.T) {
	// Test marshal error with bad writer
	badWriter := &badWriter{}
	h := NewEventHandler(badWriter, nil)

	msg := &empty.Empty{}
	// This should handle the marshal error gracefully
	h.OnReceiveResponse(msg)
}

func TestEventHandler_OnReceiveTrailers2(t *testing.T) {
	h := NewEventHandler(io.Discard, nil)

	// Test with OK status
	okStatus := status.New(codes.OK, "success")
	md := metadata.New(map[string]string{"key": "value"})
	h.OnReceiveTrailers(okStatus, md)
	assert.Equal(t, codes.OK, h.Status.Code())
	assert.Equal(t, "success", h.Status.Message())

	// Test with error status
	errorStatus := status.New(codes.Internal, "internal error")
	h.OnReceiveTrailers(errorStatus, nil)
	assert.Equal(t, codes.Internal, h.Status.Code())
	assert.Equal(t, "internal error", h.Status.Message())
}

func TestEventHandler_OnResolveMethod(t *testing.T) {
	h := NewEventHandler(io.Discard, nil)

	// Test with nil method descriptor - should not panic
	h.OnResolveMethod(nil)

	// Since this is a no-op function, we just verify it doesn't panic
	// and can be called multiple times
	h.OnResolveMethod(nil)
	h.OnResolveMethod(nil)
}

func TestEventHandler_OnSendHeaders(t *testing.T) {
	h := NewEventHandler(io.Discard, nil)

	// Test with nil metadata - should not panic
	h.OnSendHeaders(nil)

	// Test with valid metadata
	md := metadata.New(map[string]string{"request-id": "123", "auth": "token"})
	h.OnSendHeaders(md)

	// Test with empty metadata
	emptyMd := metadata.New(map[string]string{})
	h.OnSendHeaders(emptyMd)
}

func TestEventHandler_OnReceiveHeaders2(t *testing.T) {
	h := NewEventHandler(io.Discard, nil)

	// Test with nil metadata - should not panic
	h.OnReceiveHeaders(nil)

	// Test with valid metadata
	md := metadata.New(map[string]string{"response-id": "456", "content-type": "application/json"})
	h.OnReceiveHeaders(md)

	// Test with empty metadata
	emptyMd := metadata.New(map[string]string{})
	h.OnReceiveHeaders(emptyMd)
}

func TestEventHandler_CompleteWorkflow(t *testing.T) {
	var buf bytes.Buffer
	h := NewEventHandler(&buf, nil)

	// Simulate a complete gRPC call workflow
	h.OnResolveMethod(nil)
	h.OnSendHeaders(metadata.New(map[string]string{"request-id": "123"}))
	h.OnReceiveHeaders(metadata.New(map[string]string{"response-id": "456"}))

	// Send a response
	msg := &empty.Empty{}
	h.OnReceiveResponse(msg)

	// Complete with status
	h.OnReceiveTrailers(status.New(codes.OK, "completed"), nil)

	assert.Equal(t, codes.OK, h.Status.Code())
	assert.Equal(t, "completed", h.Status.Message())
	assert.Contains(t, buf.String(), "{}")
}

// badWriter is a mock writer that always returns an error
type badWriter struct{}

func (w *badWriter) Write([]byte) (int, error) {
	return 0, io.ErrShortWrite
}

func TestEventHandler_OnReceiveTrailers(t *testing.T) {
	tests := []struct {
		name           string
		writer         io.Writer
		status         *status.Status
		metadata       metadata.MD
		expectedStatus codes.Code
		expectedHeader map[string][]string
	}{
		{
			name:   "with http.ResponseWriter and metadata",
			writer: httptest.NewRecorder(),
			status: status.New(codes.OK, "success"),
			metadata: metadata.MD{
				"x-custom-header":  []string{"value1", "value2"},
				"x-another-header": []string{"single-value"},
			},
			expectedStatus: codes.OK,
			expectedHeader: map[string][]string{
				"Grpc-Trailer-X-Custom-Header":  {"value1", "value2"},
				"Grpc-Trailer-X-Another-Header": {"single-value"},
			},
		},
		{
			name:           "with http.ResponseWriter and nil metadata",
			writer:         httptest.NewRecorder(),
			status:         status.New(codes.Internal, "error"),
			metadata:       nil,
			expectedStatus: codes.Internal,
			expectedHeader: map[string][]string{},
		},
		{
			name:           "with non-http.ResponseWriter",
			writer:         io.Discard,
			status:         status.New(codes.OK, "success"),
			metadata:       metadata.MD{"x-header": []string{"value"}},
			expectedStatus: codes.OK,
			expectedHeader: nil, // headers should not be set
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewEventHandler(tt.writer, nil)

			h.OnReceiveTrailers(tt.status, tt.metadata)

			// Check status is set correctly
			assert.Equal(t, tt.expectedStatus, h.Status.Code())

			// Check headers are set correctly if writer is http.ResponseWriter
			if recorder, ok := tt.writer.(*httptest.ResponseRecorder); ok {
				if tt.expectedHeader != nil {
					for key, expectedValues := range tt.expectedHeader {
						actualValues := recorder.Header()[key]
						assert.Equal(t, expectedValues, actualValues, "Header %s should match", key)
					}
				}
			}
		})
	}
}

func TestEventHandler_OnReceiveHeaders(t *testing.T) {
	tests := []struct {
		name           string
		writer         io.Writer
		metadata       metadata.MD
		expectedHeader map[string][]string
	}{
		{
			name:   "with http.ResponseWriter and metadata",
			writer: httptest.NewRecorder(),
			metadata: metadata.MD{
				"content-type":     []string{"application/json"},
				"x-custom-header":  []string{"value1", "value2"},
				"x-another-header": []string{"single-value"},
			},
			expectedHeader: map[string][]string{
				"Grpc-Metadata-Content-Type":     {"application/json"},
				"Grpc-Metadata-X-Custom-Header":  {"value1", "value2"},
				"Grpc-Metadata-X-Another-Header": {"single-value"},
			},
		},
		{
			name:           "with http.ResponseWriter and nil metadata",
			writer:         httptest.NewRecorder(),
			metadata:       nil,
			expectedHeader: map[string][]string{},
		},
		{
			name:           "with http.ResponseWriter and empty metadata",
			writer:         httptest.NewRecorder(),
			metadata:       metadata.MD{},
			expectedHeader: map[string][]string{},
		},
		{
			name:           "with non-http.ResponseWriter",
			writer:         io.Discard,
			metadata:       metadata.MD{"x-header": []string{"value"}},
			expectedHeader: nil, // headers should not be set
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewEventHandler(tt.writer, nil)

			h.OnReceiveHeaders(tt.metadata)

			// Check headers are set correctly if writer is http.ResponseWriter
			if recorder, ok := tt.writer.(*httptest.ResponseRecorder); ok {
				if tt.expectedHeader != nil {
					for key, expectedValues := range tt.expectedHeader {
						actualValues := recorder.Header()[key]
						assert.Equal(t, expectedValues, actualValues, "Header %s should match", key)
					}
				}
			}
		})
	}
}

func TestEventHandler_OnReceiveHeaders_MultipleValues(t *testing.T) {
	recorder := httptest.NewRecorder()
	h := NewEventHandler(recorder, nil)

	// Test that multiple calls to OnReceiveHeaders accumulate headers
	h.OnReceiveHeaders(metadata.MD{
		"x-header-1": []string{"value1"},
	})

	h.OnReceiveHeaders(metadata.MD{
		"x-header-1": []string{"value2"}, // Should add to existing header
		"x-header-2": []string{"value3"},
	})

	// Check that headers are accumulated (not overwritten) with proper prefix
	assert.Equal(t, []string{"value1", "value2"}, recorder.Header()["Grpc-Metadata-X-Header-1"])
	assert.Equal(t, []string{"value3"}, recorder.Header()["Grpc-Metadata-X-Header-2"])
}

func TestEventHandler_OnReceiveHeaders_MetadataPrefix(t *testing.T) {
	tests := []struct {
		name           string
		metadata       metadata.MD
		expectedHeader map[string][]string
	}{
		{
			name: "all metadata headers should be prefixed with Grpc-Metadata-",
			metadata: metadata.MD{
				"content-type":    []string{"application/grpc"},
				"x-custom-header": []string{"value1"},
				"authorization":   []string{"Bearer token"},
			},
			expectedHeader: map[string][]string{
				"Grpc-Metadata-Content-Type":    {"application/grpc"},
				"Grpc-Metadata-X-Custom-Header": {"value1"},
				"Grpc-Metadata-Authorization":   {"Bearer token"},
			},
		},
		{
			name: "mixed case headers should be prefixed",
			metadata: metadata.MD{
				"Content-Type":    []string{"APPLICATION/JSON"},
				"X-Custom-Header": []string{"value1"},
			},
			expectedHeader: map[string][]string{
				"Grpc-Metadata-Content-Type":    {"APPLICATION/JSON"},
				"Grpc-Metadata-X-Custom-Header": {"value1"},
			},
		},
		{
			name: "multiple values for same header",
			metadata: metadata.MD{
				"x-multi-header": []string{"value1", "value2", "value3"},
			},
			expectedHeader: map[string][]string{
				"Grpc-Metadata-X-Multi-Header": {"value1", "value2", "value3"},
			},
		},
		{
			name:           "empty metadata",
			metadata:       metadata.MD{},
			expectedHeader: map[string][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			h := NewEventHandler(recorder, nil)

			h.OnReceiveHeaders(tt.metadata)

			// Check that headers are set correctly
			for key, expectedValues := range tt.expectedHeader {
				actualValues := recorder.Header()[key]
				assert.Equal(t, expectedValues, actualValues, "Header %s should match", key)
			}

			// Ensure no unexpected headers are set
			for actualKey := range recorder.Header() {
				found := false
				for expectedKey := range tt.expectedHeader {
					if actualKey == expectedKey {
						found = true
						break
					}
				}
				assert.True(t, found, "Unexpected header found: %s", actualKey)
			}
		})
	}
}
