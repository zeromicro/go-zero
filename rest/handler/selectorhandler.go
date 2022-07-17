package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/selector"
)

// SelectorHandler returns a handler that can get the selector and colors.
func SelectorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		colors := request.Header.Values("Colors")
		if len(colors) != 0 {
			ctx = selector.NewColorsContext(ctx, colors...)
		}

		selectorVal := request.Header.Values("Selector")
		if len(selectorVal) != 0 {
			ctx = selector.NewSelectorContext(ctx, selectorVal[0])
		}

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
