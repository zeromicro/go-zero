package handler

import (
	"net/http"

	"go.opentelemetry.io/otel/propagation"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/metainfo"
)

// CustomKeysHandler returns a middleware that extract custom keys from request metadata
// and inject it into request context and logger fields.
func CustomKeysHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		// try to extract custom keys from request metadata
		ctx = metainfo.CustomKeysMapPropagator.Extract(ctx, propagation.HeaderCarrier(request.Header))

		// try to inject custom keys to logger
		if info := metainfo.GetMapFromContext(ctx); len(info) > 0 {
			ctx = logx.ContextWithFields(ctx, logx.Field("custom_keys", info))
		}

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
