package internal

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestWithDialOption(t *testing.T) {
	var options ClientOptions
	agent := grpc.WithUserAgent("chrome")
	opt := WithDialOption(agent)
	opt(&options)
	assert.Contains(t, options.DialOptions, agent)
}

func TestWithTimeout(t *testing.T) {
	var options ClientOptions
	opt := WithTimeout(time.Second)
	opt(&options)
	assert.Equal(t, time.Second, options.Timeout)
}

func TestWithUnaryClientInterceptor(t *testing.T) {
	var options ClientOptions
	opt := WithUnaryClientInterceptor(func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return nil
	})
	opt(&options)
	assert.Equal(t, 1, len(options.DialOptions))
}

func TestBuildDialOptions(t *testing.T) {
	var c client
	agent := grpc.WithUserAgent("chrome")
	opts := c.buildDialOptions(WithDialOption(agent))
	assert.Contains(t, opts, agent)
}
