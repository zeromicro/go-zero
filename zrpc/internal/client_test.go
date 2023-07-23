package internal

import (
	"context"
	"net"
	"strings"
	"sync"
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

func TestWithNonBlock(t *testing.T) {
	var options ClientOptions
	opt := WithNonBlock()
	opt(&options)
	assert.True(t, options.NonBlock)
}

func TestWithStreamClientInterceptor(t *testing.T) {
	var options ClientOptions
	opt := WithStreamClientInterceptor(func(ctx context.Context, desc *grpc.StreamDesc,
		cc *grpc.ClientConn, method string, streamer grpc.Streamer,
		opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, nil
	})
	opt(&options)
	assert.Equal(t, 1, len(options.DialOptions))
}

func TestWithTransportCredentials(t *testing.T) {
	var options ClientOptions
	opt := WithTransportCredentials(nil)
	opt(&options)
	assert.Equal(t, 1, len(options.DialOptions))
}

func TestWithUnaryClientInterceptor(t *testing.T) {
	var options ClientOptions
	opt := WithUnaryClientInterceptor(func(ctx context.Context, method string, req, reply any,
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return nil
	})
	opt(&options)
	assert.Equal(t, 1, len(options.DialOptions))
}

func TestBuildDialOptions(t *testing.T) {
	c := client{
		middlewares: ClientMiddlewaresConf{
			Trace:      true,
			Duration:   true,
			Prometheus: true,
			Breaker:    true,
			Timeout:    true,
		},
	}
	agent := grpc.WithUserAgent("chrome")
	opts := c.buildDialOptions(WithDialOption(agent))
	assert.Contains(t, opts, agent)
}

func TestClientDial(t *testing.T) {
	var addr string
	var wg sync.WaitGroup
	wg.Add(1)
	server := grpc.NewServer()

	go func() {
		lis, err := net.Listen("tcp", "localhost:0")
		assert.NoError(t, err)
		defer lis.Close()
		addr = lis.Addr().String()
		wg.Done()
		server.Serve(lis)
	}()

	wg.Wait()
	c, err := NewClient(addr, ClientMiddlewaresConf{
		Trace:      true,
		Duration:   true,
		Prometheus: true,
		Breaker:    true,
		Timeout:    true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, c.Conn())
	server.Stop()
}

func TestClientDialFail(t *testing.T) {
	_, err := NewClient("localhost:54321", ClientMiddlewaresConf{
		Trace:      true,
		Duration:   true,
		Prometheus: true,
		Breaker:    true,
		Timeout:    true,
	})
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "localhost:54321"))

	_, err = NewClient("localhost:54321/fail", ClientMiddlewaresConf{
		Trace:      true,
		Duration:   true,
		Prometheus: true,
		Breaker:    true,
		Timeout:    true,
	})
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "localhost:54321/fail"))
}
