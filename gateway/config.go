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

	// Upstream is the configuration for upstream.
	Upstream struct {
		// Grpc is the target of upstream.
		Grpc zrpc.RpcClientConf
		// ProtoSet is the file of proto set, like hello.pb
		ProtoSet string `json:",optional"`
		Mapping  []struct {
			// Path is the HTTP path.
			Path string
			// Method is the gRPC method, with format of package.service/method
			Method string
		}
	}
)
