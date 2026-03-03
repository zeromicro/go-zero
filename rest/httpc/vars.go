package httpc

import "errors"

const (
	pathKey   = "path"
	formKey   = "form"
	headerKey = "header"
	jsonKey   = "json"
	slash     = "/"
	colon     = ':'
)

var (
	// ErrGetWithBody indicates that GET request with body.
	ErrGetWithBody = errors.New("HTTP GET should not have body")
	// ErrHeadWithBody indicates that HEAD request with body.
	ErrHeadWithBody = errors.New("HTTP HEAD should not have body")
)
