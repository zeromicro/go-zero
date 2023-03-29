package cors

import (
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/rest/internal/response"
)

const (
	allowOrigin      = "Access-Control-Allow-Origin"
	allOrigins       = "*"
	allowMethods     = "Access-Control-Allow-Methods"
	allowHeaders     = "Access-Control-Allow-Headers"
	allowCredentials = "Access-Control-Allow-Credentials"
	exposeHeaders    = "Access-Control-Expose-Headers"
	requestMethod    = "Access-Control-Request-Method"
	requestHeaders   = "Access-Control-Request-Headers"
	allowHeadersVal  = "Content-Type, Origin, X-CSRF-Token, Authorization, AccessToken, Token, Range"
	exposeHeadersVal = "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers"
	methods          = "GET, HEAD, POST, PATCH, PUT, DELETE"
	allowTrue        = "true"
	maxAgeHeader     = "Access-Control-Max-Age"
	maxAgeHeaderVal  = "86400"
	varyHeader       = "Vary"
	originHeader     = "Origin"
)

// NotAllowedHandler handles cross domain not allowed requests.
// At most one origin can be specified, other origins are ignored if given, default to be *.
func NotAllowedHandler(fn func(w http.ResponseWriter), origins ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gw := response.NewHeaderOnceResponseWriter(w)
		checkAndSetHeaders(gw, r, origins)
		if fn != nil {
			fn(gw)
		}

		if r.Method == http.MethodOptions {
			gw.WriteHeader(http.StatusNoContent)
		} else {
			gw.WriteHeader(http.StatusNotFound)
		}
	})
}

// Middleware returns a middleware that adds CORS headers to the response.
func Middleware(fn func(w http.Header), origins ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			checkAndSetHeaders(w, r, origins)
			if fn != nil {
				fn(w.Header())
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
			} else {
				next(w, r)
			}
		}
	}
}

func checkAndSetHeaders(w http.ResponseWriter, r *http.Request, origins []string) {
	setVaryHeaders(w, r)

	if len(origins) == 0 {
		setHeader(w, allOrigins)
		return
	}

	origin := r.Header.Get(originHeader)
	if isOriginAllowed(origins, origin) {
		setHeader(w, origin)
	}
}

func isOriginAllowed(allows []string, origin string) bool {
	origin = strings.ToLower(origin)

	for _, allow := range allows {
		if allow == allOrigins {
			return true
		}

		allow = strings.ToLower(allow)
		if origin == allow {
			return true
		}

		if strings.HasSuffix(origin, "."+allow) {
			return true
		}
	}

	return false
}

func setHeader(w http.ResponseWriter, origin string) {
	header := w.Header()
	header.Set(allowOrigin, origin)
	header.Set(allowMethods, methods)
	header.Set(allowHeaders, allowHeadersVal)
	header.Set(exposeHeaders, exposeHeadersVal)
	if origin != allOrigins {
		header.Set(allowCredentials, allowTrue)
	}
	header.Set(maxAgeHeader, maxAgeHeaderVal)
}

func setVaryHeaders(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Add(varyHeader, originHeader)
	if r.Method == http.MethodOptions {
		header.Add(varyHeader, requestMethod)
		header.Add(varyHeader, requestHeaders)
	}
}
