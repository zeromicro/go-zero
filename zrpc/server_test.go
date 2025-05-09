package zrpc

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc/internal"
	"github.com/zeromicro/go-zero/zrpc/internal/serverinterceptors"
	"google.golang.org/grpc"
)

func TestServer(t *testing.T) {
	DontLogContentForMethod("foo")
	SetServerSlowThreshold(time.Second)
	svr := MustNewServer(RpcServerConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				ServiceName: "foo",
				Mode:        "console",
			},
		},
		ListenOn:      "localhost:0",
		Etcd:          discov.EtcdConf{},
		Auth:          false,
		Redis:         redis.RedisKeyConf{},
		StrictControl: false,
		Timeout:       0,
		CpuThreshold:  0,
		Middlewares: ServerMiddlewaresConf{
			Trace:      true,
			Recover:    true,
			Stat:       true,
			Prometheus: true,
			Breaker:    true,
		},
		MethodTimeouts: []MethodTimeoutConf{
			{
				FullMethod: "/foo",
				Timeout:    time.Second,
			},
		},
	}, func(server *grpc.Server) {
	})
	svr.AddOptions(grpc.ConnectionTimeout(time.Hour))
	svr.AddUnaryInterceptors(serverinterceptors.UnaryRecoverInterceptor)
	svr.AddStreamInterceptors(serverinterceptors.StreamRecoverInterceptor)
	go svr.Start()
	svr.Stop()
}

func TestServerError(t *testing.T) {
	_, err := NewServer(RpcServerConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				ServiceName: "foo",
				Mode:        "console",
			},
		},
		ListenOn: "localhost:0",
		Etcd: discov.EtcdConf{
			Hosts: []string{"localhost"},
		},
		Auth:  true,
		Redis: redis.RedisKeyConf{},
		Middlewares: ServerMiddlewaresConf{
			Trace:      true,
			Recover:    true,
			Stat:       true,
			Prometheus: true,
			Breaker:    true,
		},
		MethodTimeouts: []MethodTimeoutConf{},
	}, func(server *grpc.Server) {
	})
	assert.NotNil(t, err)
}

func TestServer_HasEtcd(t *testing.T) {
	svr := MustNewServer(RpcServerConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				ServiceName: "foo",
				Mode:        "console",
			},
		},
		ListenOn: "localhost:0",
		Etcd: discov.EtcdConf{
			Hosts: []string{"notexist"},
			Key:   "any",
		},
		Redis: redis.RedisKeyConf{},
		Middlewares: ServerMiddlewaresConf{
			Trace:      true,
			Recover:    true,
			Stat:       true,
			Prometheus: true,
			Breaker:    true,
		},
		MethodTimeouts: []MethodTimeoutConf{},
	}, func(server *grpc.Server) {
	})
	svr.AddOptions(grpc.ConnectionTimeout(time.Hour))
	svr.AddUnaryInterceptors(serverinterceptors.UnaryRecoverInterceptor)
	svr.AddStreamInterceptors(serverinterceptors.StreamRecoverInterceptor)
	go svr.Start()
	svr.Stop()
}

func TestServer_StartFailed(t *testing.T) {
	svr := MustNewServer(RpcServerConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				ServiceName: "foo",
				Mode:        "console",
			},
		},
		ListenOn: "localhost:aaa",
		Middlewares: ServerMiddlewaresConf{
			Trace:      true,
			Recover:    true,
			Stat:       true,
			Prometheus: true,
			Breaker:    true,
		},
	}, func(server *grpc.Server) {
	})

	assert.Panics(t, svr.Start)
}

type mockedServer struct {
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
}

func (m *mockedServer) AddOptions(_ ...grpc.ServerOption) {
}

func (m *mockedServer) AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	m.streamInterceptors = append(m.streamInterceptors, interceptors...)
}

func (m *mockedServer) AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	m.unaryInterceptors = append(m.unaryInterceptors, interceptors...)
}

func (m *mockedServer) SetName(_ string) {
}

func (m *mockedServer) Start(_ internal.RegisterFn) error {
	return nil
}

func Test_setupUnaryInterceptors(t *testing.T) {
	tests := []struct {
		name string
		r    *mockedServer
		conf RpcServerConf
		len  int
	}{
		{
			name: "empty",
			r:    &mockedServer{},
			len:  0,
		},
		{
			name: "custom",
			r: &mockedServer{
				unaryInterceptors: []grpc.UnaryServerInterceptor{
					func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
						handler grpc.UnaryHandler) (interface{}, error) {
						return nil, nil
					},
				},
			},
			len: 1,
		},
		{
			name: "middleware",
			r:    &mockedServer{},
			conf: RpcServerConf{
				Middlewares: ServerMiddlewaresConf{
					Trace:      true,
					Recover:    true,
					Stat:       true,
					Prometheus: true,
					Breaker:    true,
				},
			},
			len: 5,
		},
		{
			name: "internal middleware",
			r:    &mockedServer{},
			conf: RpcServerConf{
				CpuThreshold: 900,
				Timeout:      100,
				Middlewares: ServerMiddlewaresConf{
					Trace:      true,
					Recover:    true,
					Stat:       true,
					Prometheus: true,
					Breaker:    true,
				},
			},
			len: 7,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metrics := stat.NewMetrics("abc")
			setupUnaryInterceptors(test.r, test.conf, metrics)
			assert.Equal(t, test.len, len(test.r.unaryInterceptors))
		})
	}
}

func Test_setupStreamInterceptors(t *testing.T) {
	tests := []struct {
		name string
		r    *mockedServer
		conf RpcServerConf
		len  int
	}{
		{
			name: "empty",
			r:    &mockedServer{},
			len:  0,
		},
		{
			name: "custom",
			r: &mockedServer{
				streamInterceptors: []grpc.StreamServerInterceptor{
					func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
						return handler(srv, ss)
					},
				},
			},
			len: 1,
		},
		{
			name: "middleware",
			r:    &mockedServer{},
			conf: RpcServerConf{
				Middlewares: ServerMiddlewaresConf{
					Trace:      true,
					Recover:    true,
					Stat:       true,
					Prometheus: true,
					Breaker:    true,
				},
			},
			len: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setupStreamInterceptors(test.r, test.conf)
			assert.Equal(t, test.len, len(test.r.streamInterceptors))
		})
	}
}

func Test_setupAuthInterceptors(t *testing.T) {
	t.Run("no need set auth", func(t *testing.T) {
		s := &mockedServer{}
		err := setupAuthInterceptors(s, RpcServerConf{
			Auth:  false,
			Redis: redis.RedisKeyConf{},
		})
		assert.NoError(t, err)
	})

	t.Run("redis error", func(t *testing.T) {
		s := &mockedServer{}
		err := setupAuthInterceptors(s, RpcServerConf{
			Auth:  true,
			Redis: redis.RedisKeyConf{},
		})
		assert.Error(t, err)
	})

	t.Run("works", func(t *testing.T) {
		rds := miniredis.RunT(t)
		s := &mockedServer{}
		err := setupAuthInterceptors(s, RpcServerConf{
			Auth: true,
			Redis: redis.RedisKeyConf{
				RedisConf: redis.RedisConf{
					Host: rds.Addr(),
					Type: redis.NodeType,
				},
				Key: "foo",
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(s.unaryInterceptors))
		assert.Equal(t, 1, len(s.streamInterceptors))
	})
}
