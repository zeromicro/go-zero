package gateway

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
)

type (
	Upstream struct {
		Target  string
		Mapping []struct {
			Path   string
			Method string
		}
	}

	GatewayConf struct {
		rest.RestConf
		Upstreams []Upstream
		Timeout   time.Duration `json:",default=5s"`
	}
)
