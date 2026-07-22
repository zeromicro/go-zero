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
	return strings.TrimPrefix(target.URL.Path, slashSeparator)
}

// GetHosts returns the comma-separated etcd hosts from the target URL.
// It supports two formats:
//   - New format (etcd:///h1:port,h2:port?key=k): hosts are in the URL path (empty authority)
//   - Legacy format (etcd://h1:port/key): host is in the URL authority
func GetHosts(target resolver.Target) string {
	if target.URL.Host == "" {
		// New format: hosts encoded in URL path to avoid RFC 3986 authority issues
		return GetEndpoints(target)
	}
	// Legacy format: single host in authority
	return target.URL.Host
}

// GetKey returns the etcd key from the target URL.
// It supports two formats:
//   - New format (etcd:///h1:port,h2:port?key=k): key is in the "key" query parameter
//   - Legacy format (etcd://h1:port/key): key is in the URL path
func GetKey(target resolver.Target) string {
	if target.URL.Host == "" {
		// New format: key is in the query parameter
		return target.URL.Query().Get("key")
	}
	// Legacy format: key is in the path
	return strings.Trim(target.URL.Path, slashSeparator)
}
