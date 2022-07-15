package gateway

import (
	"context"
	"net/http"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

type Server struct {
	svr       *rest.Server
	upstreams []Upstream
}

func MustNewServer(c GatewayConf) *Server {
	return &Server{
		svr:       rest.MustNewServer(c.RestConf),
		upstreams: c.Upstreams,
	}
}

func (s *Server) Start() {
	logx.Must(s.build())
	s.svr.Start()
}

func (s *Server) Stop() {
	s.svr.Stop()
}

func (s *Server) build() error {
	for _, upstream := range s.upstreams {
		zcli, err := zrpc.NewClientWithTarget(upstream.Target)
		if err != nil {
			return err
		}

		cli := grpc_reflection_v1alpha.NewServerReflectionClient(zcli.Conn())
		client := grpcreflect.NewClient(context.Background(), cli)
		source := grpcurl.DescriptorSourceFromServer(context.Background(), client)
		resolver := grpcurl.AnyResolverFromDescriptorSource(source)
		unmarshaler := jsonpb.Unmarshaler{AnyResolver: resolver, AllowUnknownFields: true}
		for _, mapping := range upstream.Mapping {
			s.svr.AddRoute(rest.Route{
				Method: http.MethodPost,
				Path:   mapping.Path,
				Handler: func(w http.ResponseWriter, r *http.Request) {
					handler := &grpcurl.DefaultEventHandler{
						Out:       w,
						Formatter: grpcurl.NewJSONFormatter(true, grpcurl.AnyResolverFromDescriptorSource(source)),
					}
					rp := grpcurl.NewJSONRequestParserWithUnmarshaler(r.Body, unmarshaler)
					grpcurl.InvokeRPC(context.Background(), source, zcli.Conn(), mapping.Method, nil, handler, rp.Next)
				},
			})
		}
	}

	return nil
}
