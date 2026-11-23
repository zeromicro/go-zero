package gateway

import (
	"net/http"
	"strings"
)

// hopByHopHeaders are the connection-specific header fields that only apply to
// a single connection, so a proxy must not forward them, see
// https://www.rfc-editor.org/rfc/rfc9110#section-7.6.1
var hopByHopHeaders = []string{
	"Connection",
	"Proxy-Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"TE",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

// copyResponseHeaders copies the headers of the upstream response to the
// client, leaving out the connection-specific ones. The client only talks to
// the gateway, so "Keep-Alive" or "Proxy-Authenticate" of the upstream would
// describe a connection it does not have.
func copyResponseHeaders(dst, src http.Header) {
	for key, values := range src {
		if isHopByHopHeader(key) {
			continue
		}

		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func isHopByHopHeader(key string) bool {
	for _, hopByHop := range hopByHopHeaders {
		if strings.EqualFold(key, hopByHop) {
			return true
		}
	}

	return false
}
