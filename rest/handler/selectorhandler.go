package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/selector"
)

func SelectorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		colors := request.Header.Values("Colors")
		if len(colors) == 0 {
			next.ServeHTTP(writer, request)
			return
		}

		ctx := request.Context()
		ctx = selector.NewColorContext(ctx, colors...)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
