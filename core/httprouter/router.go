package httprouter

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidMethod = errors.New("not a valid http method")
	ErrInvalidPath   = errors.New("path must begin with '/'")
)

type (
	Route struct {
		Path    string
		Handler http.HandlerFunc
	}

	Router interface {
		http.Handler
		Handle(method string, path string, handler http.Handler) error
		SetNotFoundHandler(handler http.Handler)
	}
)
