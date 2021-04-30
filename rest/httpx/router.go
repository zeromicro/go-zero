package httpx

import "net/http"

// Router interface represents a http router that handles http requests.
type Router interface {
	http.Handler
	Handle(method string, path string, handler http.Handler) error
	SetNotFoundHandler(handler http.Handler)
	SetNotAllowedHandler(handler http.Handler)
}
