package metainfo

import (
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestGrpcHeaderCarrier_Get(t *testing.T) {
	md := metadata.MD{
		"key1": []string{"value1"},
		"key2": []string{"value2", "value3"},
	}
	carrier := GrpcHeaderCarrier(md)

	tests := []struct {
		key      string
		expected string
	}{
		{"key1", "value1"},
		{"key2", "value2"},
		{"key3", ""},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			value := carrier.Get(tt.key)
			if value != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, value)
			}
		})
	}
}

func TestGrpcHeaderCarrier_Set(t *testing.T) {
	md := metadata.MD{}
	carrier := GrpcHeaderCarrier(md)

	carrier.Set("key1", "value1")
	carrier.Set("key2", "value2")

	expected := metadata.MD{
		"key1": []string{"value1"},
		"key2": []string{"value2"},
	}

	if !equalMetadata(metadata.MD(carrier), expected) {
		t.Errorf("expected %v, got %v", expected, carrier)
	}
}

func TestGrpcHeaderCarrier_Keys(t *testing.T) {
	md := metadata.MD{
		"key1": []string{"value1"},
		"key2": []string{"value2"},
	}
	carrier := GrpcHeaderCarrier(md)

	keys := carrier.Keys()

	expectedKeys := []string{"key1", "key2"}
	if !equalStringSlices(keys, expectedKeys) {
		t.Errorf("expected %v, got %v", expectedKeys, keys)
	}
}

// Helper function to compare two metadata.MD objects
func equalMetadata(a, b metadata.MD) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if bv, ok := b[k]; !ok || !equalStringSlices(v, bv) {
			return false
		}
	}

	return true
}

// Helper function to compare two slices of strings
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
