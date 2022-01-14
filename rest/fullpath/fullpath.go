package fullpath

import (
	"context"
	"net/http"
)

type contextKey string

var fullpathVars = contextKey("fullpathCtx")

// FullPath returns a matched route full path from given r.
// If fullpath is not found, returns empty string.
func FullPath(r *http.Request) string {
	vars, ok := r.Context().Value(fullpathVars).(string)
	if ok {
		return vars
	}
	return ""
}

// WithFullPath writes fullpath into given r and returns a new http.Request.
func WithFullPath(r *http.Request, fullPath string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), fullpathVars, fullPath))
}
