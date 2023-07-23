package clientinterceptors

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestDurationInterceptor(t *testing.T) {
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
			err := DurationInterceptor(context.Background(), "/foo", nil, nil, cc,
				func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
					opts ...grpc.CallOption) error {
					return test.err
				})
			assert.Equal(t, test.err, err)
		})
	}

	DontLogContentForMethod("/foo")
	t.Cleanup(func() {
		notLoggingContentMethods.Delete("/foo")
	})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cc := new(grpc.ClientConn)
			err := DurationInterceptor(context.Background(), "/foo", nil, nil, cc,
				func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
					opts ...grpc.CallOption) error {
					return test.err
				})
			assert.Equal(t, test.err, err)
		})
	}
}

func TestDurationInterceptorWithSlowThreshold(t *testing.T) {
	SetSlowThreshold(time.Microsecond)
	t.Cleanup(func() {
		SetSlowThreshold(defaultSlowThreshold)
	})

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
			err := DurationInterceptor(context.Background(), "/foo", nil, nil, cc,
				func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
					opts ...grpc.CallOption) error {
					time.Sleep(time.Millisecond * 10)
					return test.err
				})
			assert.Equal(t, test.err, err)
		})
	}

	DontLogContentForMethod("/foo")
	t.Cleanup(func() {
		notLoggingContentMethods.Delete("/foo")
	})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cc := new(grpc.ClientConn)
			err := DurationInterceptor(context.Background(), "/foo", nil, nil, cc,
				func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
					opts ...grpc.CallOption) error {
					time.Sleep(time.Millisecond * 10)
					return test.err
				})
			assert.Equal(t, test.err, err)
		})
	}
}

func TestSetSlowThreshold(t *testing.T) {
	assert.Equal(t, defaultSlowThreshold, slowThreshold.Load())
	SetSlowThreshold(time.Second)
	assert.Equal(t, time.Second, slowThreshold.Load())
}
