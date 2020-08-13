package httpx

import "net/http"

const (
	xForwardFor = "X-Forward-For"
	xRealIp     = "X-Real-IP"
)

type ClientIP struct {
	RemoteIP  string
	ForwardIP string
	RealIP    string
}

// Returns the peer address, supports X-Forward-For
func GetRemoteAddr(r *http.Request) string {
	v := r.Header.Get(xForwardFor)
	if len(v) > 0 {
		return v
	}
	return r.RemoteAddr
}

// Returns the client address, supports X-Forward-For,X-Real-IP
func GetClientIps(r *http.Request) *ClientIP {
	remoteIp, _, _ := net.SplitHostPort(r.RemoteAddr)
	xForwardFor := r.Header.Get(xForwardFor)
	xRealIp := r.Header.Get(xRealIp)
	return &ClientIP{remoteIp, xForwardFor, xRealIp}
}
