package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/md"
)

// MdHandler returns a handler that can get the selector and colors.
func MdHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := md.Extract(request.Context(), md.HeaderCarrier(request.Header))
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
