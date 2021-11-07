package cors

import "net/http"

const (
	allowOrigin      = "Access-Control-Allow-Origin"
	allOrigins       = "*"
	allowMethods     = "Access-Control-Allow-Methods"
	allowHeaders     = "Access-Control-Allow-Headers"
	allowCredentials = "Access-Control-Allow-Credentials"
	exposeHeaders    = "Access-Control-Expose-Headers"
	allowHeadersVal  = "Content-Type, Origin, X-CSRF-Token, Authorization, AccessToken, Token, Range"
	exposeHeadersVal = "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers"
	methods          = "GET, HEAD, POST, PATCH, PUT, DELETE"
	allowTrue        = "true"
	maxAgeHeader     = "Access-Control-Max-Age"
	maxAgeHeaderVal  = "86400"
)

// Handler handles cross domain not allowed requests.
// At most one origin can be specified, other origins are ignored if given, default to be *.
func Handler(origin ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setHeader(w, getOrigin(origin))

		if r.Method != http.MethodOptions {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	})
}

// Middleware returns a middleware that adds CORS headers to the response.
func Middleware(origin ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			setHeader(w, getOrigin(origin))

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
			} else {
				next(w, r)
			}
		}
	}
}

func getOrigin(origins []string) string {
	if len(origins) > 0 {
		return origins[0]
	} else {
		return allOrigins
	}
}

func setHeader(w http.ResponseWriter, origin string) {
	w.Header().Set(allowOrigin, origin)
	w.Header().Set(allowMethods, methods)
	w.Header().Set(allowHeaders, allowHeadersVal)
	w.Header().Set(exposeHeaders, exposeHeadersVal)
	w.Header().Set(allowCredentials, allowTrue)
	w.Header().Set(maxAgeHeader, maxAgeHeaderVal)
}
