package httpx

import "net/http"

type Router interface {
	http.Handler
	Handle(method string, path string, handler http.Handler) error
	SetNotFoundHandler(handler http.Handler)
}
