package gateway

import (
	"errors"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

var errMissingMethodPath = errors.New("gateway: RouteMapping is missing Method or Path, " +
	"set top-level fields or use Match block")

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

	// Match defines the matching conditions for a route.
	// Currently supports Path and Method matching.
	Match struct {
		// Path is the HTTP path pattern for matching.
		Path string `json:",optional"`
		// Method is the HTTP method for matching, like GET, POST, PUT, DELETE.
		Method string `json:",optional"`
	}

	// RouteMapping is a mapping between a gateway route and an upstream rpc method.
	RouteMapping struct {
		// Method is the HTTP method, like GET, POST, PUT, DELETE.
		Method string `json:",optional"`
		// Path is the HTTP path.
		Path string `json:",optional"`
		// RpcPath is the gRPC rpc method, with format of package.service/method, optional.
		// If the mapping is for HTTP, it's not necessary.
		RpcPath string `json:",optional"`
		// Match defines the matching conditions for the route.
		// If Match is set, Match.Path and Match.Method take precedence over
		// the top-level Path and Method fields.
		Match *Match `json:",optional"`
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

// GetMethod returns the resolved HTTP method for the route mapping.
// If Match is set and Match.Method is non-empty, it takes precedence.
func (m RouteMapping) GetMethod() string {
	if m.Match != nil && len(m.Match.Method) > 0 {
		return m.Match.Method
	}
	return m.Method
}

// GetPath returns the resolved HTTP path for the route mapping.
// If Match is set and Match.Path is non-empty, it takes precedence.
func (m RouteMapping) GetPath() string {
	if m.Match != nil && len(m.Match.Path) > 0 {
		return m.Match.Path
	}
	return m.Path
}

// Validate checks that the route mapping has a non-empty Method and Path,
// resolved from either the Match block or the top-level fields.
func (m RouteMapping) Validate() error {
	if len(m.GetMethod()) == 0 || len(m.GetPath()) == 0 {
		return errMissingMethodPath
	}
	return nil
}
