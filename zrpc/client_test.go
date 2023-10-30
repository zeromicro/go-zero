package zrpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/internal/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func init() {
	logx.Disable()
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	mock.RegisterDepositServiceServer(server, &mock.DepositServer{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestDepositServer_Deposit(t *testing.T) {
	tests := []struct {
		name    string
		amount  float32
		timeout time.Duration
		res     *mock.DepositResponse
		errCode codes.Code
		errMsg  string
	}{
		{
			name:    "invalid request with negative amount",
			amount:  -1.11,
			errCode: codes.InvalidArgument,
			errMsg:  fmt.Sprintf("cannot deposit %v", -1.11),
		},
		{
			name:    "valid request with non negative amount",
			res:     &mock.DepositResponse{Ok: true},
			errCode: codes.OK,
		},
		{
			name:    "valid request with long handling time",
			amount:  2000.00,
			errCode: codes.DeadlineExceeded,
			errMsg:  "context deadline exceeded",
		},
		{
			name:    "valid request with timeout call option",
			amount:  2000.00,
			timeout: time.Second * 3,
			res:     &mock.DepositResponse{Ok: true},
			errCode: codes.OK,
			errMsg:  "",
		},
	}

	directClient := MustNewClient(
		RpcClientConf{
			Endpoints: []string{"foo"},
			App:       "foo",
			Token:     "bar",
			Timeout:   1000,
			Middlewares: ClientMiddlewaresConf{
				Trace:      true,
				Duration:   true,
				Prometheus: true,
				Breaker:    true,
				Timeout:    true,
			},
		},
		WithDialOption(grpc.WithContextDialer(dialer())),
		WithUnaryClientInterceptor(func(ctx context.Context, method string, req, reply any,
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)
	nonBlockClient := MustNewClient(
		RpcClientConf{
			Endpoints: []string{"foo"},
			App:       "foo",
			Token:     "bar",
			Timeout:   1000,
			NonBlock:  true,
			Middlewares: ClientMiddlewaresConf{
				Trace:      true,
				Duration:   true,
				Prometheus: true,
				Breaker:    true,
				Timeout:    true,
			},
		},
		WithDialOption(grpc.WithContextDialer(dialer())),
		WithUnaryClientInterceptor(func(ctx context.Context, method string, req, reply any,
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)
	tarConfClient := MustNewClient(
		RpcClientConf{
			Target:        "foo",
			App:           "foo",
			Token:         "bar",
			Timeout:       1000,
			KeepaliveTime: time.Second * 15,
			Middlewares: ClientMiddlewaresConf{
				Trace:      true,
				Duration:   true,
				Prometheus: true,
				Breaker:    true,
				Timeout:    true,
			},
		},
		WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		WithDialOption(grpc.WithContextDialer(dialer())),
		WithUnaryClientInterceptor(func(ctx context.Context, method string, req, reply any,
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)
	targetClient, err := NewClientWithTarget("foo",
		WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		WithDialOption(grpc.WithContextDialer(dialer())), WithUnaryClientInterceptor(
			func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
				invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				return invoker(ctx, method, req, reply, cc, opts...)
			}), WithTimeout(1000*time.Millisecond))
	assert.Nil(t, err)
	clients := []Client{
		directClient,
		nonBlockClient,
		tarConfClient,
		targetClient,
	}
	DontLogClientContentForMethod("foo")
	SetClientSlowThreshold(time.Second)

	for _, tt := range tests {
		tt := tt
		for _, client := range clients {
			client := client
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				cli := mock.NewDepositServiceClient(client.Conn())
				request := &mock.DepositRequest{Amount: tt.amount}

				var (
					ctx      = context.Background()
					response *mock.DepositResponse
					err      error
				)

				if tt.timeout > 0 {
					response, err = cli.Deposit(ctx, request, WithCallTimeout(tt.timeout))
				} else {
					response, err = cli.Deposit(ctx, request)
				}

				if response != nil {
					assert.True(t, len(response.String()) > 0)
					if response.GetOk() != tt.res.GetOk() {
						t.Error("response: expected", tt.res.GetOk(), "received", response.GetOk())
					}
				}
				if err != nil {
					if e, ok := status.FromError(err); ok {
						if e.Code() != tt.errCode {
							t.Error("error code: expected", codes.InvalidArgument, "received", e.Code())
						}
						if e.Message() != tt.errMsg {
							t.Error("error message: expected", tt.errMsg, "received", e.Message())
						}
					}
				}
			})
		}
	}
}

func TestNewClientWithError(t *testing.T) {
	_, err := NewClient(
		RpcClientConf{
			App:     "foo",
			Token:   "bar",
			Timeout: 1000,
		},
		WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		WithDialOption(grpc.WithContextDialer(dialer())),
		WithUnaryClientInterceptor(func(ctx context.Context, method string, req, reply any,
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)
	assert.NotNil(t, err)

	_, err = NewClient(
		RpcClientConf{
			Etcd: discov.EtcdConf{
				Hosts: []string{"localhost:2379"},
				Key:   "mock",
			},
			App:     "foo",
			Token:   "bar",
			Timeout: 1,
		},
		WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		WithDialOption(grpc.WithContextDialer(dialer())),
		WithUnaryClientInterceptor(func(ctx context.Context, method string, req, reply any,
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)
	assert.NotNil(t, err)
}

func TestNewClientWithTarget(t *testing.T) {
	_, err := NewClientWithTarget("",
		WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		WithDialOption(grpc.WithContextDialer(dialer())),
		WithUnaryClientInterceptor(func(ctx context.Context, method string, req, reply any,
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}))

	assert.NotNil(t, err)
}
