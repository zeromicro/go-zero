package handler

import (
	"context"
	"net/http"
)

// I18nHandler returns a middleware that recovers if panic happens.
func I18nHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// add lang to context
		ctx = context.WithValue(ctx, "lang", r.Header.Get("Accept-Language"))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
