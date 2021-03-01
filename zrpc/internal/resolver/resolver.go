package resolver

import (
	"fmt"

	"google.golang.org/grpc/resolver"
)

const (
	// DirectScheme stands for direct schema.
	DirectScheme = "direct"
	// DiscovSchema stands for discov schema.
	DiscovScheme = "discov"
	// EnpointSepChar is the separator cha in endpoints.
	EndpointSepChar = ','

	subsetSize = 32
)

var (
	// EnpointSep is the separator string in endpoints.
	EndpointSep = fmt.Sprintf("%c", EndpointSepChar)

	dirBuilder directBuilder
	disBuilder discovBuilder
)

// RegisterResolver registers the direct and discov schemas to the resolver.
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
