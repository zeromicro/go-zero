package resolver

import (
	"fmt"

	"google.golang.org/grpc/resolver"
)

const (
	DirectScheme    = "direct"
	DiscovScheme    = "discov"
	DiscovK8sScheme = "k8s"
	EndpointSepChar = ','
	subsetSize      = 32
)

var (
	EndpointSep   = fmt.Sprintf("%c", EndpointSepChar)
	dirBuilder    directBuilder
	disBuilder    discovBuilder
	disK8sBuilder discovK8sBuilder
)

func RegisterResolver() {
	resolver.Register(&dirBuilder)
	resolver.Register(&disBuilder)
	resolver.Register(&disK8sBuilder)
}

type nopResolver struct {
	cc resolver.ClientConn
}

func (r *nopResolver) Close() {
}

func (r *nopResolver) ResolveNow(options resolver.ResolveNowOptions) {
}
