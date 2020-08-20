package internal

import (
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/rpcx/internal/resolver"
)

func BuildDirectTarget(endpoints []string) string {
	return fmt.Sprintf("%s:///%s", resolver.DirectScheme, strings.Join(
		endpoints, fmt.Sprint(resolver.EndpointSep)))
}

func BuildDiscovTarget(endpoints []string, key string) string {
	return fmt.Sprintf("%s://%s/%s", resolver.DiscovScheme, strings.Join(
		endpoints, fmt.Sprint(resolver.EndpointSep)), key)
}
