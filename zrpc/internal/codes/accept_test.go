package codes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAccept(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		accept bool
	}{
		{
			name:   "nil error",
			err:    nil,
			accept: true,
		},
		{
			name:   "deadline error",
			err:    status.Error(codes.DeadlineExceeded, "deadline"),
			accept: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.accept, Acceptable(test.err))
		})
	}
}
