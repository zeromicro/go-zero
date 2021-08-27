package codes

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAccept(t *testing.T) {
	err := errors.New("custom error")
	AddBreakerErrors(err)
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
		}, {
			name:   "custom error",
			err:    err,
			accept: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.accept, Acceptable(test.err))
		})
	}
}
