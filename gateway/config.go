package gateway

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
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
		// Target is the target of upstream, like etcd://localhost:2379/hello.rpc
		Target  string
		Mapping []struct {
			// Path is the HTTP path.
			Path string
			// Method is the gRPC method, with format of package.service/method
			Method string
		}
	}
)
