package metainfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			assert.Equal(t, tt.expected, carrier.Get(tt.key))
		})
	}
}

func TestGrpcHeaderCarrier_Set(t *testing.T) {
	md := metadata.MD{}
	carrier := GrpcHeaderCarrier(md)

	carrier.Set("key1", "value1")
	carrier.Set("key2", "value2")

	assert.Equal(t, metadata.MD{
		"key1": []string{"value1"},
		"key2": []string{"value2"},
	}, md)
}

func TestGrpcHeaderCarrier_Keys(t *testing.T) {
	md := metadata.MD{
		"key1": []string{"value1"},
		"key2": []string{"value2"},
	}
	carrier := GrpcHeaderCarrier(md)

	assert.ElementsMatch(t, []string{"key1", "key2"}, carrier.Keys())
}
