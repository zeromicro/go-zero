package gateway

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/gateway/internal"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc/codes"
)

type (
	// Server is a gateway server.
	Server struct {
		*rest.Server
		upstreams     []Upstream
		conns         []zrpc.Client
		processHeader func(http.Header) []string
		dialer        func(conf zrpc.RpcClientConf) zrpc.Client
	}

	// Option defines the method to customize Server.
	Option func(svr *Server)
)

// MustNewServer creates a new gateway server.
func MustNewServer(c GatewayConf, opts ...Option) *Server {
	svr := &Server{
		upstreams: c.Upstreams,
		Server:    rest.MustNewServer(c.RestConf),
	}
	for _, opt := range opts {
		opt(svr)
	}

	return svr
}

// Start starts the gateway server.
func (s *Server) Start() {
	logx.Must(s.build())
	s.Server.Start()
}

// Stop stops the gateway server.
// To get a graceful shutdown, it stops the HTTP server first, then closes gRPC connections.
func (s *Server) Stop() {
	// stop the HTTP server first, then close gRPC connections.
	// in case the gRPC server is stopped first,
	// the HTTP server may still be running to accept requests.
	s.Server.Stop()

	group := threading.NewRoutineGroup()
	for _, conn := range s.conns {
		// new variable to avoid closure problems, can be removed after go 1.22
		// see https://golang.org/doc/faq#closures_and_goroutines
		conn := conn
		group.Run(func() {
			// ignore the error when closing the connection
			_ = conn.Conn().Close()
		})
	}
	group.Wait()
}

func (s *Server) build() error {
	if err := s.ensureUpstreamNames(); err != nil {
		return err
	}

	return mr.MapReduceVoid(func(source chan<- Upstream) {
		for _, up := range s.upstreams {
			source <- up
		}
	}, func(up Upstream, writer mr.Writer[rest.Route], cancel func(error)) {
		var cli zrpc.Client
		if s.dialer != nil {
			cli = s.dialer(up.Grpc)
		} else {
			cli = zrpc.MustNewClient(up.Grpc)
		}
		s.conns = append(s.conns, cli)

		source, err := s.createDescriptorSource(cli, up)
		if err != nil {
			cancel(fmt.Errorf("%s: %w", up.Name, err))
			return
		}

		methods, err := internal.GetMethods(source)
		if err != nil {
			cancel(fmt.Errorf("%s: %w", up.Name, err))
			return
		}

		resolver := grpcurl.AnyResolverFromDescriptorSource(source)
		for _, m := range methods {
			if len(m.HttpMethod) > 0 && len(m.HttpPath) > 0 {
				writer.Write(rest.Route{
					Method:  m.HttpMethod,
					Path:    m.HttpPath,
					Handler: s.buildHandler(source, resolver, cli, m.RpcPath),
				})
			}
		}

		methodSet := make(map[string]struct{})
		for _, m := range methods {
			methodSet[m.RpcPath] = struct{}{}
		}
		for _, m := range up.Mappings {
			if _, ok := methodSet[m.RpcPath]; !ok {
				cancel(fmt.Errorf("%s: rpc method %s not found", up.Name, m.RpcPath))
				return
			}

			writer.Write(rest.Route{
				Method:  strings.ToUpper(m.Method),
				Path:    m.Path,
				Handler: s.buildHandler(source, resolver, cli, m.RpcPath),
			})
		}
	}, func(pipe <-chan rest.Route, cancel func(error)) {
		for route := range pipe {
			s.Server.AddRoute(route)
		}
	})
}

func (s *Server) buildHandler(source grpcurl.DescriptorSource, resolver jsonpb.AnyResolver,
	cli zrpc.Client, rpcPath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		parser, err := internal.NewRequestParser(r, resolver)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		w.Header().Set(httpx.ContentType, httpx.JsonContentType)
		handler := internal.NewEventHandler(w, resolver)
		if err := grpcurl.InvokeRPC(r.Context(), source, cli.Conn(), rpcPath, s.prepareMetadata(r.Header),
			handler, parser.Next); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}

		st := handler.Status
		if st.Code() != codes.OK {
			httpx.ErrorCtx(r.Context(), w, st.Err())
		}
	}
}

func (s *Server) createDescriptorSource(cli zrpc.Client, up Upstream) (grpcurl.DescriptorSource, error) {
	var source grpcurl.DescriptorSource
	var err error

	if len(up.ProtoSets) > 0 {
		source, err = grpcurl.DescriptorSourceFromProtoSets(up.ProtoSets...)
		if err != nil {
			return nil, err
		}
	} else {
		client := grpcreflect.NewClientAuto(context.Background(), cli.Conn())
		source = grpcurl.DescriptorSourceFromServer(context.Background(), client)
	}

	return source, nil
}

func (s *Server) ensureUpstreamNames() error {
	for i := 0; i < len(s.upstreams); i++ {
		target, err := s.upstreams[i].Grpc.BuildTarget()
		if err != nil {
			return err
		}

		s.upstreams[i].Name = target
	}

	return nil
}

func (s *Server) prepareMetadata(header http.Header) []string {
	vals := internal.ProcessHeaders(header)
	if s.processHeader != nil {
		vals = append(vals, s.processHeader(header)...)
	}

	return vals
}

// WithHeaderProcessor sets a processor to process request headers.
// The returned headers are used as metadata to invoke the RPC.
func WithHeaderProcessor(processHeader func(http.Header) []string) func(*Server) {
	return func(s *Server) {
		s.processHeader = processHeader
	}
}

// withDialer sets a dialer to create a gRPC client.
func withDialer(dialer func(conf zrpc.RpcClientConf) zrpc.Client) func(*Server) {
	return func(s *Server) {
		s.dialer = dialer
	}
}
