package gateway

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/gateway/internal"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpc"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc/codes"
)

const defaultHttpScheme = "http"

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
		// up.Grpc and up.Http are exclusive
		if up.Grpc != nil {
			s.buildGrpcRoute(up, writer, cancel)
		} else if up.Http != nil {
			s.buildHttpRoute(up, writer)
		}
	}, func(pipe <-chan rest.Route, cancel func(error)) {
		for route := range pipe {
			s.Server.AddRoute(route)
		}
	})
}

func (s *Server) buildGrpcHandler(source grpcurl.DescriptorSource, resolver jsonpb.AnyResolver,
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

func (s *Server) buildGrpcRoute(up Upstream, writer mr.Writer[rest.Route], cancel func(error)) {
	var cli zrpc.Client
	if s.dialer != nil {
		cli = s.dialer(*up.Grpc)
	} else {
		cli = zrpc.MustNewClient(*up.Grpc)
	}
	s.conns = append(s.conns, cli)

	source, err := createDescriptorSource(cli, up)
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
				Handler: s.buildGrpcHandler(source, resolver, cli, m.RpcPath),
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
			Handler: s.buildGrpcHandler(source, resolver, cli, m.RpcPath),
		})
	}
}

func (s *Server) buildHttpHandler(target *HttpClientConf) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(httpx.ContentType, httpx.JsonContentType)
		req, err := buildRequestWithNewTarget(r, target)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// set the timeout if it's configured, take effect only if it's greater than 0
		// and less than the deadline of the original request
		if target.Timeout > 0 {
			timeout := time.Duration(target.Timeout) * time.Millisecond
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			req = req.WithContext(ctx)
		}

		resp, err := httpc.DoRequest(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		defer resp.Body.Close()

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(resp.StatusCode)
		if _, err = io.Copy(w, resp.Body); err != nil {
			// log the error with original request info
			logc.Error(r.Context(), err)
		}
	}
}

func (s *Server) buildHttpRoute(up Upstream, writer mr.Writer[rest.Route]) {
	for _, m := range up.Mappings {
		writer.Write(rest.Route{
			Method:  strings.ToUpper(m.Method),
			Path:    m.Path,
			Handler: s.buildHttpHandler(up.Http),
		})
	}
}

func (s *Server) ensureUpstreamNames() error {
	for i := 0; i < len(s.upstreams); i++ {
		if len(s.upstreams[i].Name) > 0 {
			continue
		}

		if s.upstreams[i].Grpc != nil {
			target, err := s.upstreams[i].Grpc.BuildTarget()
			if err != nil {
				return err
			}

			s.upstreams[i].Name = target
		} else if s.upstreams[i].Http != nil {
			s.upstreams[i].Name = s.upstreams[i].Http.Target
		}
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

func buildRequestWithNewTarget(r *http.Request, target *HttpClientConf) (*http.Request, error) {
	u := *r.URL
	u.Host = target.Target
	if len(u.Scheme) == 0 {
		u.Scheme = defaultHttpScheme
	}

	if len(target.Prefix) > 0 {
		var err error
		u.Path, err = url.JoinPath(target.Prefix, u.Path)
		if err != nil {
			return nil, err
		}
	}

	newReq := &http.Request{
		Method:        r.Method,
		URL:           &u,
		Header:        r.Header.Clone(),
		Proto:         r.Proto,
		ProtoMajor:    r.ProtoMajor,
		ProtoMinor:    r.ProtoMinor,
		ContentLength: r.ContentLength,
		Body:          io.NopCloser(r.Body),
	}

	// make sure the context is passed to the new request
	return newReq.WithContext(r.Context()), nil
}

func createDescriptorSource(cli zrpc.Client, up Upstream) (grpcurl.DescriptorSource, error) {
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

// withDialer sets a dialer to create a gRPC client.
func withDialer(dialer func(conf zrpc.RpcClientConf) zrpc.Client) func(*Server) {
	return func(s *Server) {
		s.dialer = dialer
	}
}
