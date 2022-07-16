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
		Upstreams []upstream
		Timeout   time.Duration `json:",default=5s"`
	}

	// mapping is a mapping between a gateway route and a upstream rpc method.
	mapping struct {
		// Method is the HTTP method, like GET, POST, PUT, DELETE.
		Method string
		// Path is the HTTP path.
		Path string
		// Rpc is the gRPC rpc method, with format of package.service/method
		Rpc string
	}
	// upstream is the configuration for upstream.
	upstream struct {
		// Grpc is the target of upstream.
		Grpc zrpc.RpcClientConf
		// ProtoSet is the file of proto set, like hello.pb
		ProtoSet string `json:",optional"`
		Mapping  []mapping
	}
)
