package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestWithStreamClientInterceptors(t *testing.T) {
	opts := WithStreamClientInterceptors()
	assert.NotNil(t, opts)
}

func TestWithUnaryClientInterceptors(t *testing.T) {
	opts := WithUnaryClientInterceptors()
	assert.NotNil(t, opts)
}

func TestChainStreamClientInterceptors_zero(t *testing.T) {
	var vals []int
	interceptors := chainStreamClientInterceptors()
	_, err := interceptors(context.Background(), nil, new(grpc.ClientConn), "/foo",
		func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
			opts ...grpc.CallOption) (grpc.ClientStream, error) {
			vals = append(vals, 1)
			return nil, nil
		})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{1}, vals)
}

func TestChainStreamClientInterceptors_one(t *testing.T) {
	var vals []int
	interceptors := chainStreamClientInterceptors(func(ctx context.Context, desc *grpc.StreamDesc,
		cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (
		grpc.ClientStream, error) {
		vals = append(vals, 1)
		return streamer(ctx, desc, cc, method, opts...)
	})
	_, err := interceptors(context.Background(), nil, new(grpc.ClientConn), "/foo",
		func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
			opts ...grpc.CallOption) (grpc.ClientStream, error) {
			vals = append(vals, 2)
			return nil, nil
		})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{1, 2}, vals)
}

func TestChainStreamClientInterceptors_more(t *testing.T) {
	var vals []int
	interceptors := chainStreamClientInterceptors(func(ctx context.Context, desc *grpc.StreamDesc,
		cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (
		grpc.ClientStream, error) {
		vals = append(vals, 1)
		return streamer(ctx, desc, cc, method, opts...)
	}, func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
		streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		vals = append(vals, 2)
		return streamer(ctx, desc, cc, method, opts...)
	})
	_, err := interceptors(context.Background(), nil, new(grpc.ClientConn), "/foo",
		func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
			opts ...grpc.CallOption) (grpc.ClientStream, error) {
			vals = append(vals, 3)
			return nil, nil
		})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{1, 2, 3}, vals)
}

func TestWithUnaryClientInterceptors_zero(t *testing.T) {
	var vals []int
	interceptors := chainUnaryClientInterceptors()
	err := interceptors(context.Background(), "/foo", nil, nil, new(grpc.ClientConn),
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
			opts ...grpc.CallOption) error {
			vals = append(vals, 1)
			return nil
		})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{1}, vals)
}

func TestWithUnaryClientInterceptors_one(t *testing.T) {
	var vals []int
	interceptors := chainUnaryClientInterceptors(func(ctx context.Context, method string, req,
		reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		vals = append(vals, 1)
		return invoker(ctx, method, req, reply, cc, opts...)
	})
	err := interceptors(context.Background(), "/foo", nil, nil, new(grpc.ClientConn),
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
			opts ...grpc.CallOption) error {
			vals = append(vals, 2)
			return nil
		})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{1, 2}, vals)
}

func TestWithUnaryClientInterceptors_more(t *testing.T) {
	var vals []int
	interceptors := chainUnaryClientInterceptors(func(ctx context.Context, method string, req,
		reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		vals = append(vals, 1)
		return invoker(ctx, method, req, reply, cc, opts...)
	}, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		vals = append(vals, 2)
		return invoker(ctx, method, req, reply, cc, opts...)
	})
	err := interceptors(context.Background(), "/foo", nil, nil, new(grpc.ClientConn),
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
			opts ...grpc.CallOption) error {
			vals = append(vals, 3)
			return nil
		})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{1, 2, 3}, vals)
}
