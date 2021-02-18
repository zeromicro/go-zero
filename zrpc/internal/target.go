package internal

import (
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/zrpc/internal/resolver"
)

func BuildDirectTarget(endpoints []string) string {
	return fmt.Sprintf("%s:///%s", resolver.DirectScheme,
		strings.Join(endpoints, resolver.EndpointSep))
}

func BuildDiscovTarget(endpoints []string, key string) string {
	return fmt.Sprintf("%s://%s/%s", resolver.DiscovScheme,
		strings.Join(endpoints, resolver.EndpointSep), key)
}

// The target format is: kubernetes://service-name.namespace:8080/
func BuildDiscovk8sTarget(name string, namespace string, port uint32) string {
	return fmt.Sprintf("%s://%s.%s:%v/", resolver.DiscovK8sScheme, name, namespace, port)
}
