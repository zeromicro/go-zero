package clientinterceptors

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/breaker"
	"github.com/tal-tech/go-zero/core/stat"
	rcodes "github.com/tal-tech/go-zero/zrpc/internal/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	stat.SetReporter(nil)
}

type mockError struct {
	st *status.Status
}

func (m mockError) GRPCStatus() *status.Status {
	return m.st
}

func (m mockError) Error() string {
	return "mocked error"
}

func TestBreakerInterceptorNotFound(t *testing.T) {
	err := mockError{st: status.New(codes.NotFound, "any")}
	for i := 0; i < 1000; i++ {
		assert.Equal(t, err, breaker.DoWithAcceptable("call", func() error {
			return err
		}, rcodes.Acceptable))
	}
}

func TestBreakerInterceptorDeadlineExceeded(t *testing.T) {
	err := mockError{st: status.New(codes.DeadlineExceeded, "any")}
	errs := make(map[error]int)
	for i := 0; i < 1000; i++ {
		e := breaker.DoWithAcceptable("call", func() error {
			return err
		}, rcodes.Acceptable)
		errs[e]++
	}
	assert.Equal(t, 2, len(errs))
	assert.True(t, errs[err] > 0)
	assert.True(t, errs[breaker.ErrServiceUnavailable] > 0)
}

func TestBreakerInterceptor(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "nil",
			err:  nil,
		},
		{
			name: "with error",
			err:  errors.New("mock"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cc := new(grpc.ClientConn)
			err := BreakerInterceptor(context.Background(), "/foo", nil, nil, cc,
				func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
					opts ...grpc.CallOption) error {
					return test.err
				})
			assert.Equal(t, test.err, err)
		})
	}
}
