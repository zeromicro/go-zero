package context

import (
	"context"
	"net/http"
)

var pathVars = contextKey("pathVars")

func Vars(r *http.Request) map[string]string {
	vars, ok := r.Context().Value(pathVars).(map[string]string)
	if ok {
		return vars
	}

	return nil
}

func WithPathVars(r *http.Request, params map[string]string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), pathVars, params))
}

type contextKey string

func (c contextKey) String() string {
	return "rest/internal/context context key" + string(c)
}
