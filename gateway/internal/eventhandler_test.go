package internal

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestEventHandler(t *testing.T) {
	h := NewEventHandler(io.Discard, nil)
	h.OnResolveMethod(nil)
	h.OnSendHeaders(nil)
	h.OnReceiveHeaders(nil)
	h.OnReceiveTrailers(status.New(codes.OK, ""), nil)
	assert.Equal(t, codes.OK, h.Status.Code())
	h.OnReceiveResponse(nil)
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
				"X-Custom-Header":  {"value1", "value2"},
				"X-Another-Header": {"single-value"},
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
				"Content-Type":     {"application/json"},
				"X-Custom-Header":  {"value1", "value2"},
				"X-Another-Header": {"single-value"},
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
	
	// Check that headers are accumulated (not overwritten)
	assert.Equal(t, []string{"value1", "value2"}, recorder.Header()["X-Header-1"])
	assert.Equal(t, []string{"value3"}, recorder.Header()["X-Header-2"])
}
