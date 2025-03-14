package gateway

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type (
	// GatewayConf is the configuration for gateway.
	GatewayConf struct {
		rest.RestConf
		Upstreams []Upstream
	}

	// HttpClientConf is the configuration for an HTTP client.
	HttpClientConf struct {
		Target  string
		Prefix  string `json:",optional"`
		Timeout int64  `json:",default=3000"`
	}

	// RouteMapping is a mapping between a gateway route and an upstream rpc method.
	RouteMapping struct {
		// Method is the HTTP method, like GET, POST, PUT, DELETE.
		Method string
		// Path is the HTTP path.
		Path string
		// RpcPath is the gRPC rpc method, with format of package.service/method, optional.
		// If the mapping is for HTTP, it's not necessary.
		RpcPath string `json:",optional"`
	}

	// Upstream is the configuration for an upstream.
	Upstream struct {
		// Name is the name of the upstream.
		Name string `json:",optional"`
		// Grpc is the target of the upstream.
		Grpc *zrpc.RpcClientConf `json:",optional"`
		// Http is the target of the upstream.
		Http *HttpClientConf `json:",optional=!grpc"`
		// ProtoSets is the file list of proto set, like [hello.pb].
		// if your proto file import another proto file, you need to write multi-file slice,
		// like [hello.pb, common.pb].
		ProtoSets []string `json:",optional"`
		// Mappings is the mapping between gateway routes and Upstream methods.
		// Keep it blank if annotations are added in rpc methods.
		Mappings []RouteMapping `json:",optional"`
	}
)
