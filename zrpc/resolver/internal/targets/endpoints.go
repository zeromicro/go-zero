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
