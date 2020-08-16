package httpx

import "net/http"

const xForwardFor = "X-Forward-For"

// Returns the peer address, supports X-Forward-For
func GetRemoteAddr(r *http.Request) string {
	v := r.Header.Get(xForwardFor)
	if len(v) > 0 {
		return v
	}
	return r.RemoteAddr
}
