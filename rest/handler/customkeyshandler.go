package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/metainfo"
	"go.opentelemetry.io/otel/propagation"
)

// CustomKeysHandler returns a middleware that extract custom keys from request metadata
// and inject it into request context and logger fields.
func CustomKeysHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// extract custom keys from request metadata
		ctx = metainfo.CustomKeysMapPropagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

		// inject custom keys to logger
		if info := metainfo.GetMapFromContext(ctx); len(info) > 0 {
			ctx = logx.ContextWithFields(ctx, logx.Field(metainfo.LogKey, info))
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
