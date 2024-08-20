//go:build !no_k8s

package internal

import "google.golang.org/grpc/resolver"

var k8sResolverBuilder kubeBuilder

// RegisterResolver registers the direct, etcd, discov and k8s schemes to the resolver.
func RegisterResolver() {
	register()
	resolver.Register(&k8sResolverBuilder)
}
