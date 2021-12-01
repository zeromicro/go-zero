package internal

import (
	"fmt"

	"google.golang.org/grpc/resolver"
)

const (
	// DirectScheme stands for direct scheme.
	DirectScheme = "direct"
	// DiscovScheme stands for discov scheme.
	DiscovScheme = "discov"
	// EtcdScheme stands for etcd scheme.
	EtcdScheme = "etcd"
	// KubernetesScheme stands for k8s scheme.
	KubernetesScheme = "k8s"
	// EndpointSepChar is the separator cha in endpoints.
	EndpointSepChar = ','

	subsetSize = 32
)

var (
	// EndpointSep is the separator string in endpoints.
	EndpointSep = fmt.Sprintf("%c", EndpointSepChar)

	directResolverBuilder directBuilder
	discovResolverBuilder discovBuilder
	etcdResolverBuilder   etcdBuilder
	k8sResolverBuilder    kubeBuilder
)

// RegisterResolver registers the direct and discov schemes to the resolver.
func RegisterResolver() {
	resolver.Register(&directResolverBuilder)
	resolver.Register(&discovResolverBuilder)
	resolver.Register(&etcdResolverBuilder)
	resolver.Register(&k8sResolverBuilder)
}

type nopResolver struct {
	cc resolver.ClientConn
}

func (r *nopResolver) Close() {
}

func (r *nopResolver) ResolveNow(options resolver.ResolveNowOptions) {
}
