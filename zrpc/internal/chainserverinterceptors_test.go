package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestWithStreamServerInterceptors(t *testing.T) {
	opts := WithStreamServerInterceptors()
	assert.NotNil(t, opts)
}

func TestWithUnaryServerInterceptors(t *testing.T) {
	opts := WithUnaryServerInterceptors()
	assert.NotNil(t, opts)
}

func TestChainStreamServerInterceptors_zero(t *testing.T) {
	var vals []int
	interceptors := chainStreamServerInterceptors()
	err := interceptors(nil, nil, nil, func(srv interface{}, stream grpc.ServerStream) error {
		vals = append(vals, 1)
		return nil
	})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{1}, vals)
}

func TestChainStreamServerInterceptors_one(t *testing.T) {
	var vals []int
	interceptors := chainStreamServerInterceptors(func(srv interface{}, ss grpc.ServerStream,
		info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		vals = append(vals, 1)
		return handler(srv, ss)
	})
	err := interceptors(nil, nil, nil, func(srv interface{}, stream grpc.ServerStream) error {
		vals = append(vals, 2)
		return nil
	})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{1, 2}, vals)
}

func TestChainStreamServerInterceptors_more(t *testing.T) {
	var vals []int
	interceptors := chainStreamServerInterceptors(func(srv interface{}, ss grpc.ServerStream,
		info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		vals = append(vals, 1)
		return handler(srv, ss)
	}, func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		vals = append(vals, 2)
		return handler(srv, ss)
	})
	err := interceptors(nil, nil, nil, func(srv interface{}, stream grpc.ServerStream) error {
		vals = append(vals, 3)
		return nil
	})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{1, 2, 3}, vals)
}

func TestChainUnaryServerInterceptors_zero(t *testing.T) {
	var vals []int
	interceptors := chainUnaryServerInterceptors()
	_, err := interceptors(context.Background(), nil, nil,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			vals = append(vals, 1)
			return nil, nil
		})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{1}, vals)
}

func TestChainUnaryServerInterceptors_one(t *testing.T) {
	var vals []int
	interceptors := chainUnaryServerInterceptors(func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		vals = append(vals, 1)
		return handler(ctx, req)
	})
	_, err := interceptors(context.Background(), nil, nil,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			vals = append(vals, 2)
			return nil, nil
		})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{1, 2}, vals)
}

func TestChainUnaryServerInterceptors_more(t *testing.T) {
	var vals []int
	interceptors := chainUnaryServerInterceptors(func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		vals = append(vals, 1)
		return handler(ctx, req)
	}, func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		vals = append(vals, 2)
		return handler(ctx, req)
	})
	_, err := interceptors(context.Background(), nil, nil,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			vals = append(vals, 3)
			return nil, nil
		})
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{1, 2, 3}, vals)
}
