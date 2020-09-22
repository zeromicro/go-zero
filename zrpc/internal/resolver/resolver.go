package resolver

import (
	"fmt"

	"google.golang.org/grpc/resolver"
)

const (
	DirectScheme    = "direct"
	DiscovScheme    = "discov"
	EndpointSepChar = ','
	subsetSize      = 32
)

var (
	EndpointSep = fmt.Sprintf("%c", EndpointSepChar)
	dirBuilder  directBuilder
	disBuilder  discovBuilder
)

func RegisterResolver() {
	resolver.Register(&dirBuilder)
	resolver.Register(&disBuilder)
}

type nopResolver struct {
	cc resolver.ClientConn
}

func (r *nopResolver) Close() {
}

func (r *nopResolver) ResolveNow(options resolver.ResolveNowOptions) {
}
