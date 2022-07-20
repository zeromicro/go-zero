package gateway

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/gateway/internal"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

type (
	// Server is a gateway server.
	Server struct {
		*rest.Server
		upstreams     []Upstream
		timeout       time.Duration
		processHeader func(http.Header) []string
	}

	// Option defines the method to customize Server.
	Option func(svr *Server)
)

// MustNewServer creates a new gateway server.
func MustNewServer(c GatewayConf, opts ...Option) *Server {
	svr := &Server{
		Server:    rest.MustNewServer(c.RestConf),
		upstreams: c.Upstreams,
		timeout:   c.Timeout,
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
func (s *Server) Stop() {
	s.Server.Stop()
}

func (s *Server) build() error {
	return mr.MapReduceVoid(func(source chan<- interface{}) {
		for _, up := range s.upstreams {
			source <- up
		}
	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		up := item.(Upstream)
		cli := zrpc.MustNewClient(up.Grpc)
		source, err := s.createDescriptorSource(cli, up)
		if err != nil {
			cancel(err)
			return
		}

		methods, err := internal.GetMethods(source)
		if err != nil {
			cancel(err)
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
		for _, m := range up.Mapping {
			if _, ok := methodSet[m.RpcPath]; !ok {
				cancel(fmt.Errorf("rpc method %s not found", m.RpcPath))
				return
			}

			writer.Write(rest.Route{
				Method:  strings.ToUpper(m.Method),
				Path:    m.Path,
				Handler: s.buildHandler(source, resolver, cli, m.RpcPath),
			})
		}
	}, func(pipe <-chan interface{}, cancel func(error)) {
		for item := range pipe {
			route := item.(rest.Route)
			s.Server.AddRoute(route)
		}
	})
}

func (s *Server) buildHandler(source grpcurl.DescriptorSource, resolver jsonpb.AnyResolver,
	cli zrpc.Client, rpcPath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler := &grpcurl.DefaultEventHandler{
			Out: w,
			Formatter: grpcurl.NewJSONFormatter(true,
				grpcurl.AnyResolverFromDescriptorSource(source)),
		}
		parser, err := internal.NewRequestParser(r, resolver)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		timeout := internal.GetTimeout(r.Header, s.timeout)
		ctx, can := context.WithTimeout(r.Context(), timeout)
		defer can()

		w.Header().Set(httpx.ContentType, httpx.JsonContentType)
		if err := grpcurl.InvokeRPC(ctx, source, cli.Conn(), rpcPath, s.prepareMetadata(r.Header),
			handler, parser.Next); err != nil {
			httpx.Error(w, err)
		}
	}
}

func (s *Server) createDescriptorSource(cli zrpc.Client, up Upstream) (grpcurl.DescriptorSource, error) {
	var source grpcurl.DescriptorSource
	var err error

	if len(up.ProtoSet) > 0 {
		source, err = grpcurl.DescriptorSourceFromProtoSets(up.ProtoSet)
		if err != nil {
			return nil, err
		}
	} else {
		refCli := grpc_reflection_v1alpha.NewServerReflectionClient(cli.Conn())
		client := grpcreflect.NewClient(context.Background(), refCli)
		source = grpcurl.DescriptorSourceFromServer(context.Background(), client)
	}

	return source, nil
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
