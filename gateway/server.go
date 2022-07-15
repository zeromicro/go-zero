package gateway

import (
	"context"
	"net/http"
	"time"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// Server is a gateway server.
type Server struct {
	svr       *rest.Server
	upstreams []Upstream
	timeout   time.Duration
}

// MustNewServer creates a new gateway server.
func MustNewServer(c GatewayConf) *Server {
	return &Server{
		svr:       rest.MustNewServer(c.RestConf),
		upstreams: c.Upstreams,
		timeout:   c.Timeout,
	}
}

// Start starts the gateway server.
func (s *Server) Start() {
	logx.Must(s.build())
	s.svr.Start()
}

// Stop stops the gateway server.
func (s *Server) Stop() {
	s.svr.Stop()
}

func (s *Server) build() error {
	return mr.MapReduceVoid(func(source chan<- interface{}) {
		for _, upstream := range s.upstreams {
			source <- upstream
		}
	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		upstream := item.(Upstream)
		cli := zrpc.MustNewClient(upstream.Grpc)
		source, err := s.createDescriptorSource(cli, upstream)
		if err != nil {
			cancel(err)
			return
		}

		resolver := grpcurl.AnyResolverFromDescriptorSource(source)
		unmarshaler := jsonpb.Unmarshaler{AnyResolver: resolver, AllowUnknownFields: true}
		for _, mapping := range upstream.Mapping {
			writer.Write(rest.Route{
				Method: http.MethodPost,
				Path:   mapping.Path,
				Handler: func(w http.ResponseWriter, r *http.Request) {
					handler := &grpcurl.DefaultEventHandler{
						Out: w,
						Formatter: grpcurl.NewJSONFormatter(true,
							grpcurl.AnyResolverFromDescriptorSource(source)),
					}
					rp := grpcurl.NewJSONRequestParserWithUnmarshaler(r.Body, unmarshaler)
					ctx, can := context.WithTimeout(r.Context(), s.timeout)
					defer can()
					if err := grpcurl.InvokeRPC(ctx, source, cli.Conn(), mapping.Method,
						nil, handler, rp.Next); err != nil {
						httpx.Error(w, err)
					}
				},
			})
		}
	}, func(pipe <-chan interface{}, cancel func(error)) {
		for item := range pipe {
			route := item.(rest.Route)
			s.svr.AddRoute(route)
		}
	})
}

func (s *Server) createDescriptorSource(cli zrpc.Client, upstream Upstream) (grpcurl.DescriptorSource, error) {
	var source grpcurl.DescriptorSource
	var err error

	if len(upstream.ProtoSet) > 0 {
		source, err = grpcurl.DescriptorSourceFromProtoSets(upstream.ProtoSet)
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
