package gateway

import (
	"context"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/internal/mock"
	"github.com/zeromicro/go-zero/rest/httpc"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/test/bufconn"
)

func init() {
	logx.Disable()
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	mock.RegisterDepositServiceServer(server, &mock.DepositServer{})

	reflection.Register(server)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestMustNewServer(t *testing.T) {
	var c GatewayConf
	assert.NoError(t, conf.FillDefault(&c))
	// avoid popup alert on macos for asking permissions
	c.DevServer.Host = "localhost"
	c.Host = "localhost"
	c.Port = 18881

	s := MustNewServer(c)
	s.upstreams = []*upstream{
		{
			Upstream: Upstream{
				Mappings: []RouteMapping{
					{
						Method:  "get",
						Path:    "/deposit/:amount",
						RpcPath: "mock.DepositService/Deposit",
					},
				},
			},
			client: zrpc.MustNewClient(
				zrpc.RpcClientConf{
					Endpoints: []string{"foo"},
					Timeout:   1000,
					Middlewares: zrpc.ClientMiddlewaresConf{
						Trace:      true,
						Duration:   true,
						Prometheus: true,
						Breaker:    true,
						Timeout:    true,
					},
				},
				zrpc.WithDialOption(grpc.WithContextDialer(dialer())),
			),
		},
	}

	assert.NoError(t, s.buildUpstream())
	go s.Server.Start()

	time.Sleep(time.Millisecond * 100)

	resp, err := httpc.Do(context.Background(), http.MethodGet, "http://localhost:18881/deposit/100", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = httpc.Do(context.Background(), http.MethodGet, "http://localhost:18881/deposit_fail/100", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
