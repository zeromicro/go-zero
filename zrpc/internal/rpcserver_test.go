package internal

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/internal/mock"
	"google.golang.org/grpc"
)

func TestRpcServer(t *testing.T) {
	metrics := stat.NewMetrics("foo")
	server := NewRpcServer("localhost:54321", ServerMiddlewaresConf{
		Trace:      true,
		Recover:    true,
		Stat:       true,
		Prometheus: true,
		Breaker:    true,
	}, WithMetrics(metrics), WithRpcHealth(true))
	server.SetName("mock")
	var wg, wgDone sync.WaitGroup
	var grpcServer *grpc.Server
	var lock sync.Mutex
	wg.Add(1)
	wgDone.Add(1)
	go func() {
		err := server.Start(func(server *grpc.Server) {
			lock.Lock()
			mock.RegisterDepositServiceServer(server, new(mock.DepositServer))
			grpcServer = server
			lock.Unlock()
			wg.Done()
		})
		assert.Nil(t, err)
		wgDone.Done()
	}()

	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	lock.Lock()
	grpcServer.GracefulStop()
	lock.Unlock()

	proc.WrapUp()
	wgDone.Wait()
}

func TestRpcServer_WithBadAddress(t *testing.T) {
	server := NewRpcServer("localhost:111111", ServerMiddlewaresConf{
		Trace:      true,
		Recover:    true,
		Stat:       true,
		Prometheus: true,
		Breaker:    true,
	}, WithRpcHealth(true))
	server.SetName("mock")
	err := server.Start(func(server *grpc.Server) {
		mock.RegisterDepositServiceServer(server, new(mock.DepositServer))
	})
	assert.NotNil(t, err)

	proc.WrapUp()
}

func TestRpcServer_buildUnaryInterceptor(t *testing.T) {
	tests := []struct {
		name string
		r    *rpcServer
		len  int
	}{
		{
			name: "empty",
			r: &rpcServer{
				baseRpcServer: &baseRpcServer{},
			},
			len: 0,
		},
		{
			name: "custom",
			r: &rpcServer{
				baseRpcServer: &baseRpcServer{
					unaryInterceptors: []grpc.UnaryServerInterceptor{
						func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
							handler grpc.UnaryHandler) (interface{}, error) {
							return nil, nil
						},
					},
				},
			},
			len: 1,
		},
		{
			name: "middleware",
			r: &rpcServer{
				baseRpcServer: &baseRpcServer{
					unaryInterceptors: []grpc.UnaryServerInterceptor{
						func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
							handler grpc.UnaryHandler) (interface{}, error) {
							return nil, nil
						},
					},
				},
				middlewares: ServerMiddlewaresConf{
					Trace:      true,
					Recover:    true,
					Stat:       true,
					Prometheus: true,
					Breaker:    true,
				},
			},
			len: 6,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.len, len(test.r.buildUnaryInterceptors()))
		})
	}
}

func TestRpcServer_buildStreamInterceptor(t *testing.T) {
	tests := []struct {
		name string
		r    *rpcServer
		len  int
	}{
		{
			name: "empty",
			r: &rpcServer{
				baseRpcServer: &baseRpcServer{},
			},
			len: 0,
		},
		{
			name: "custom",
			r: &rpcServer{
				baseRpcServer: &baseRpcServer{
					streamInterceptors: []grpc.StreamServerInterceptor{
						func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo,
							handler grpc.StreamHandler) error {
							return nil
						},
					},
				},
			},
			len: 1,
		},
		{
			name: "middleware",
			r: &rpcServer{
				baseRpcServer: &baseRpcServer{
					streamInterceptors: []grpc.StreamServerInterceptor{
						func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo,
							handler grpc.StreamHandler) error {
							return nil
						},
					},
				},
				middlewares: ServerMiddlewaresConf{
					Trace:   true,
					Recover: true,
					Breaker: true,
				},
			},
			len: 4,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.len, len(test.r.buildStreamInterceptors()))
		})
	}
}
