package gateway

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type (
	// GatewayConf is the configuration for gateway.
	GatewayConf struct {
		rest.RestConf
		Upstreams []Upstream
		Timeout   time.Duration `json:",default=5s"`
	}

	// RouteMapping is a mapping between a gateway route and an upstream rpc method.
	RouteMapping struct {
		// Method is the HTTP method, like GET, POST, PUT, DELETE.
		Method string
		// Path is the HTTP path.
		Path string
		// RpcPath is the gRPC rpc method, with format of package.service/method
		RpcPath string
	}

	// Upstream is the configuration for an upstream.
	Upstream struct {
		// Grpc is the target of the upstream.
		Grpc zrpc.RpcClientConf
		// ProtoSet is the file of proto set, like hello.pb
		ProtoSet string `json:",optional"`
		// Mapping is the mapping between gateway routes and Upstream rpc methods.
		// Keep it blank if annotations are added in rpc methods.
		Mapping []RouteMapping `json:",optional"`
	}
)
