package codec

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzip(t *testing.T) {
	var buf bytes.Buffer
	for i := 0; i < 10000; i++ {
		fmt.Fprint(&buf, i)
	}

	bs := Gzip(buf.Bytes())
	actual, err := Gunzip(bs)

	assert.Nil(t, err)
	assert.True(t, len(bs) < buf.Len())
	assert.Equal(t, buf.Bytes(), actual)
}

func TestGunzip(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expected    []byte
		expectedErr error
	}{
		{
			name: "valid input",
			input: func() []byte {
				var buf bytes.Buffer
				gz := gzip.NewWriter(&buf)
				gz.Write([]byte("hello"))
				gz.Close()
				return buf.Bytes()
			}(),
			expected:    []byte("hello"),
			expectedErr: nil,
		},
		{
			name:        "invalid input",
			input:       []byte("invalid input"),
			expected:    nil,
			expectedErr: gzip.ErrHeader,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Gunzip(test.input)

			if !bytes.Equal(result, test.expected) {
				t.Errorf("unexpected result: %v", result)
			}

			if !errors.Is(err, test.expectedErr) {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
