package fullpath

import (
	"context"
	"net/http"
)

type contextKey string

var fullpathVars = contextKey("fullpathCtx")

// FullPath returns fullpath from given r.
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
