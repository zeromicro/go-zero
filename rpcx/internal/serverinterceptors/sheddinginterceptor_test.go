package serverinterceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/load"
	"github.com/tal-tech/go-zero/core/stat"
	"google.golang.org/grpc"
)

func TestUnarySheddingInterceptor(t *testing.T) {
	tests := []struct {
		name      string
		allow     bool
		handleErr error
		expect    error
	}{
		{
			name:      "allow",
			allow:     true,
			handleErr: nil,
			expect:    nil,
		},
		{
			name:      "allow",
			allow:     true,
			handleErr: context.DeadlineExceeded,
			expect:    context.DeadlineExceeded,
		},
		{
			name:      "reject",
			allow:     false,
			handleErr: nil,
			expect:    load.ErrServiceOverloaded,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			shedder := mockedShedder{allow: test.allow}
			metrics := stat.NewMetrics("mock")
			interceptor := UnarySheddingInterceptor(shedder, metrics)
			_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
				FullMethod: "/",
			}, func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, test.handleErr
			})
			assert.Equal(t, test.expect, err)
		})
	}
}

type mockedShedder struct {
	allow bool
}

func (m mockedShedder) Allow() (load.Promise, error) {
	if m.allow {
		return mockedPromise{}, nil
	} else {
		return nil, load.ErrServiceOverloaded
	}
}

type mockedPromise struct {
}

func (m mockedPromise) Pass() {
}

func (m mockedPromise) Fail() {
}
