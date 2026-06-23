package resolver

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/zeromicro/go-zero/zrpc/resolver/internal"
)

// BuildDirectTarget returns a string that represents the given endpoints with direct schema.
func BuildDirectTarget(endpoints []string) string {
	return fmt.Sprintf("%s:///%s", internal.DirectScheme,
		strings.Join(endpoints, internal.EndpointSep))
}

// BuildDiscovTarget returns a string that represents the given endpoints with discov schema.
// The format is etcd:///host1:port,host2:port?key=<etcd-key> to avoid placing comma-separated
// hosts in the URI authority, which Go 1.26+ rejects per RFC 3986.
func BuildDiscovTarget(endpoints []string, key string) string {
	return fmt.Sprintf("%s:///%s?key=%s", internal.EtcdScheme,
		strings.Join(endpoints, internal.EndpointSep), url.QueryEscape(key))
}
