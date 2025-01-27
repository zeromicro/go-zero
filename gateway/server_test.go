package gateway

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/logx/logtest"
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
	// avoid popup alert on MacOS for asking permissions
	c.DevServer.Host = "localhost"
	c.Host = "localhost"
	c.Port = 18881

	s := MustNewServer(c, withDialer(func(conf zrpc.RpcClientConf) zrpc.Client {
		return zrpc.MustNewClient(conf, zrpc.WithDialOption(grpc.WithContextDialer(dialer())))
	}), WithHeaderProcessor(func(header http.Header) []string {
		return []string{"foo"}
	}))
	s.upstreams = []Upstream{
		{
			Mappings: []RouteMapping{
				{
					Method:  "get",
					Path:    "/deposit/:amount",
					RpcPath: "mock.DepositService/Deposit",
				},
			},
			Grpc: &zrpc.RpcClientConf{
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
		},
	}

	assert.NoError(t, s.build())
	go s.Server.Start()
	defer s.Stop()

	time.Sleep(time.Millisecond * 200)

	resp, err := httpc.Do(context.Background(), http.MethodGet, "http://localhost:18881/deposit/100", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = httpc.Do(context.Background(), http.MethodGet, "http://localhost:18881/deposit_fail/100", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestServer_ensureUpstreamNames(t *testing.T) {
	var s = Server{
		upstreams: []Upstream{
			{
				Grpc: &zrpc.RpcClientConf{
					Target: "target",
				},
			},
		},
	}

	assert.NoError(t, s.ensureUpstreamNames())
	assert.Equal(t, "target", s.upstreams[0].Name)
}

func TestServer_ensureUpstreamNames_badEtcd(t *testing.T) {
	var s = Server{
		upstreams: []Upstream{
			{
				Grpc: &zrpc.RpcClientConf{
					Etcd: discov.EtcdConf{},
				},
			},
		},
	}

	logtest.PanicOnFatal(t)
	assert.Panics(t, func() {
		s.Start()
	})
}

func TestHttpToHttp(t *testing.T) {
	server := startTestServer(t)
	defer server.Close()

	var c GatewayConf
	assert.NoError(t, conf.FillDefault(&c))
	c.DevServer.Host = "localhost"
	c.Host = "localhost"
	c.Port = 18882

	s := MustNewServer(c)
	s.upstreams = []Upstream{
		{
			Name: "test",
			Mappings: []RouteMapping{
				{
					Method: "get",
					Path:   "/api/ping",
				},
			},
			Http: &HttpClientConf{
				Target:  "localhost:45678",
				Timeout: 3000,
			},
		},
		{
			Mappings: []RouteMapping{
				{
					Method: "get",
					Path:   "/ping",
				},
			},
			Http: &HttpClientConf{
				Target: "localhost:45678",
				Prefix: "/api",
			},
		},
	}

	go s.Start()
	defer s.Stop()

	time.Sleep(time.Millisecond * 200)

	t.Run("/api/ping", func(t *testing.T) {
		resp, err := httpc.Do(context.Background(), http.MethodGet,
			"http://localhost:18882/api/ping", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if assert.NoError(t, err) {
			assert.Equal(t, "pong", string(body))
		}
	})

	t.Run("/ping", func(t *testing.T) {
		resp, err := httpc.Do(context.Background(), http.MethodGet,
			"http://localhost:18882/ping", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if assert.NoError(t, err) {
			assert.Equal(t, "pong", string(body))
		}
	})

	t.Run("no upstream", func(t *testing.T) {
		resp, err := httpc.Do(context.Background(), http.MethodGet,
			"http://localhost:18882/ping/bad", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestHttpToHttpBadUpstream(t *testing.T) {
	var c GatewayConf
	assert.NoError(t, conf.FillDefault(&c))
	c.DevServer.Host = "localhost"
	c.Host = "localhost"
	c.Port = 18883

	s := MustNewServer(c)
	s.upstreams = []Upstream{
		{
			Mappings: []RouteMapping{
				{
					Method: "get",
					Path:   "/api/ping",
				},
			},
			Http: &HttpClientConf{
				Target: "localhost:45678",
				Prefix: "\x7f/api",
			},
		},
	}

	go s.Start()
	defer s.Stop()

	time.Sleep(time.Millisecond * 200)

	t.Run("/api/ping", func(t *testing.T) {
		resp, err := httpc.Do(context.Background(), http.MethodGet,
			"http://localhost:18883/api/ping", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestHttpToHttpBadWriter(t *testing.T) {
	t.Run("bad url", func(t *testing.T) {
		handler := new(Server).buildHttpHandler(&HttpClientConf{
			Target:  "http://example.com",
			Timeout: 3000,
		})
		w := httptest.NewRecorder()
		handler.ServeHTTP(&badResponseWriter{w},
			httptest.NewRequest(http.MethodGet, "http://localhost:18884", nil))
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("bad url", func(t *testing.T) {
		var c GatewayConf
		assert.NoError(t, conf.FillDefault(&c))
		c.DevServer.Host = "localhost"
		c.Host = "localhost"
		c.Port = 18884

		s := MustNewServer(c)
		s.upstreams = []Upstream{
			{
				Mappings: []RouteMapping{
					{
						Method: "get",
						Path:   "/api/ping",
					},
				},
				Http: &HttpClientConf{
					Target: "localhost:45678",
					Prefix: "\x7f/api",
				},
			},
		}

		go s.Start()
		defer s.Stop()

		handler := new(Server).buildHttpHandler(&HttpClientConf{
			Target:  "localhost:18884",
			Timeout: 3000,
		})
		w := httptest.NewRecorder()
		handler.ServeHTTP(&badResponseWriter{w},
			httptest.NewRequest(http.MethodGet, "http://localhost:18884/api/ping", nil))
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// Handler function for the root route
func pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("pong"))
}

func startTestServer(t *testing.T) *http.Server {
	http.HandleFunc("/api/ping", pingHandler)

	server := &http.Server{
		Addr:    ":45678",
		Handler: http.DefaultServeMux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("failed to start server: %v", err)
		}
	}()

	return server
}

type badResponseWriter struct {
	http.ResponseWriter
}

func (w *badResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("bad writer")
}
