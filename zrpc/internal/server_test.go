package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestBaseRpcServer_AddOptions(t *testing.T) {
	server := newBaseRpcServer("foo", &rpcServerOptions{})
	var opt grpc.EmptyServerOption
	server.AddOptions(opt)
	assert.Contains(t, server.options, opt)
}

func TestBaseRpcServer_AddStreamInterceptors(t *testing.T) {
	server := newBaseRpcServer("foo", &rpcServerOptions{})
	var vals []int
	f := func(_ any, _ grpc.ServerStream, _ *grpc.StreamServerInfo, _ grpc.StreamHandler) error {
		vals = append(vals, 1)
		return nil
	}
	server.AddStreamInterceptors(f)
	for _, each := range server.streamInterceptors {
		assert.Nil(t, each(nil, nil, nil, nil))
	}
	assert.ElementsMatch(t, []int{1}, vals)
}

func TestBaseRpcServer_AddUnaryInterceptors(t *testing.T) {
	server := newBaseRpcServer("foo", &rpcServerOptions{})
	var vals []int
	f := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		resp any, err error) {
		vals = append(vals, 1)
		return nil, nil
	}
	server.AddUnaryInterceptors(f)
	for _, each := range server.unaryInterceptors {
		_, err := each(context.Background(), nil, nil, nil)
		assert.Nil(t, err)
	}
	assert.ElementsMatch(t, []int{1}, vals)
}
