package cors

import (
	"bufio"
	"errors"
	"net"
	"net/http"
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
		gw := &guardedResponseWriter{w: w}
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

type guardedResponseWriter struct {
	w           http.ResponseWriter
	wroteHeader bool
}

func (w *guardedResponseWriter) Flush() {
	if flusher, ok := w.w.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *guardedResponseWriter) Header() http.Header {
	return w.w.Header()
}

// Hijack implements the http.Hijacker interface.
// This expands the Response to fulfill http.Hijacker if the underlying http.ResponseWriter supports it.
func (w *guardedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacked, ok := w.w.(http.Hijacker); ok {
		return hijacked.Hijack()
	}

	return nil, nil, errors.New("server doesn't support hijacking")
}

func (w *guardedResponseWriter) Write(bytes []byte) (int, error) {
	return w.w.Write(bytes)
}

func (w *guardedResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}

	w.w.WriteHeader(code)
	w.wroteHeader = true
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
	for _, o := range allows {
		if o == allOrigins {
			return true
		}

		if o == origin {
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
