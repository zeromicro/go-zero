package targets

import (
	"google.golang.org/grpc/resolver"
)

const slashSeparator = '/'

// GetAuthority returns the authority of the target.
func GetAuthority(target resolver.Target) string {
	return target.URL.Host
}

// GetEndpoints returns the endpoints from the given target.
func GetEndpoints(target resolver.Target) string {
	if target.URL.Path[0] == slashSeparator {
		return target.URL.Path[1:]
	}

	return target.URL.Path
}
