package targets

import (
	"strings"

	"google.golang.org/grpc/resolver"
)

const slashSeparator = "/"

// GetAuthority returns the authority of the target.
func GetAuthority(target resolver.Target) string {
	return target.URL.Host
}

// GetEndpoints returns the endpoints from the given target.
func GetEndpoints(target resolver.Target) string {
	return strings.Trim(target.URL.Path, slashSeparator)
}

// GetHosts returns the comma-separated etcd hosts from the target URL path.
// Used for etcd/discov targets where hosts are encoded in the URL path to
// avoid RFC 3986 authority parsing issues with comma-separated hosts.
func GetHosts(target resolver.Target) string {
	return strings.Trim(target.URL.Path, slashSeparator)
}

// GetKey returns the etcd key from the "key" query parameter of the target URL.
func GetKey(target resolver.Target) string {
	return target.URL.Query().Get("key")
}
