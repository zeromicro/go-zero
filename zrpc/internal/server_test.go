package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stat"
	"google.golang.org/grpc"
)

func TestBaseRpcServer_AddOptions(t *testing.T) {
	metrics := stat.NewMetrics("foo")
	server := newBaseRpcServer("foo", metrics)
	server.SetName("bar")
	var opt grpc.EmptyServerOption
	server.AddOptions(opt)
	assert.Contains(t, server.options, opt)
}

func TestBaseRpcServer_AddStreamInterceptors(t *testing.T) {
	metrics := stat.NewMetrics("foo")
	server := newBaseRpcServer("foo", metrics)
	server.SetName("bar")
	var vals []int
	f := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
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
	metrics := stat.NewMetrics("foo")
	server := newBaseRpcServer("foo", metrics)
	server.SetName("bar")
	var vals []int
	f := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		resp interface{}, err error) {
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
