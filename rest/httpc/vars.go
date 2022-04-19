package httpc

import "errors"

const (
	pathKey         = "path"
	formKey         = "form"
	headerKey       = "header"
	jsonKey         = "json"
	slash           = "/"
	colon           = ':'
	contentType     = "Content-Type"
	applicationJson = "application/json"
)

// ErrGetWithBody indicates that GET request with body.
var ErrGetWithBody = errors.New("HTTP GET should not have body")
