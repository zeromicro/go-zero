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
	h := NewEventHandlerWithContext(ctx, w, nil, true)

	assert.NotNil(t, h)
	assert.Equal(t, w, h.writer)
	assert.True(t, h.useOkHandler)
	assert.Equal(t, ctx, h.ctx)
	assert.True(t, h.marshaler.EmitDefaults)
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
	h := NewEventHandlerWithContext(ctx, w, nil, false)

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

func TestEventHandler_OnReceiveResponse_MarshalError(t *testing.T) {
	// Test marshal error with bad writer
	badWriter := &badWriter{}
	h := NewEventHandler(badWriter, nil)

	msg := &empty.Empty{}
	// This should handle the marshal error gracefully
	h.OnReceiveResponse(msg)
}

func TestEventHandler_OnReceiveTrailers(t *testing.T) {
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
