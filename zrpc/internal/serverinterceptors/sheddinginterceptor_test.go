package serverinterceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/load"
	"github.com/zeromicro/go-zero/core/stat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnarySheddingInterceptor(t *testing.T) {
	tests := []struct {
		name       string
		allow      bool
		handleErr  error
		expect     error
		expectFail bool
	}{
		{
			name:      "allow",
			allow:     true,
			handleErr: nil,
			expect:    nil,
		},
		{
			name:       "context deadline exceeded",
			allow:      true,
			handleErr:  context.DeadlineExceeded,
			expect:     context.DeadlineExceeded,
			expectFail: true,
		},
		{
			// UnaryTimeoutInterceptor reports timeouts as a status error instead of
			// context.DeadlineExceeded, which must be counted as a failure as well.
			name:       "grpc deadline exceeded",
			allow:      true,
			handleErr:  status.Error(codes.DeadlineExceeded, context.DeadlineExceeded.Error()),
			expect:     status.Error(codes.DeadlineExceeded, context.DeadlineExceeded.Error()),
			expectFail: true,
		},
		{
			name:      "reject",
			allow:     false,
			handleErr: nil,
			expect:    status.Error(codes.ResourceExhausted, load.ErrServiceOverloaded.Error()),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			promise := &mockedPromise{}
			shedder := mockedShedder{
				allow:   test.allow,
				promise: promise,
			}
			metrics := stat.NewMetrics("mock")
			interceptor := UnarySheddingInterceptor(shedder, metrics)
			_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
				FullMethod: "/",
			}, func(ctx context.Context, req any) (any, error) {
				return nil, test.handleErr
			})
			assert.Equal(t, test.expect, err)

			if test.expectFail {
				assert.True(t, promise.failed)
				assert.False(t, promise.passed)
			} else if test.allow {
				assert.True(t, promise.passed)
				assert.False(t, promise.failed)
			}
		})
	}
}

type mockedShedder struct {
	allow   bool
	promise load.Promise
}

func (m mockedShedder) Allow() (load.Promise, error) {
	if m.allow {
		return m.promise, nil
	}

	return nil, load.ErrServiceOverloaded
}

type mockedPromise struct {
	passed bool
	failed bool
}

func (m *mockedPromise) Pass() {
	m.passed = true
}

func (m *mockedPromise) Fail() {
	m.failed = true
}
