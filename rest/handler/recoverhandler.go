package handler

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/zeromicro/go-zero/rest/internal"
	"github.com/zeromicro/go-zero/rest/internal/response"
)

// RecoverCallback defines the method of recover callback.
type RecoverCallback func(w http.ResponseWriter, r *http.Request, recoverRes interface{})

// RecoverHandler returns a middleware that recovers if panic happens.
func RecoverHandler(recoverCallback RecoverCallback) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if result := recover(); result != nil {
					internal.Error(r, fmt.Sprintf("%v\n%s", result, debug.Stack()))

					writer := response.NewHeaderOnceResponseWriter(w)
					if recoverCallback != nil {
						doRecoverCallback(
							func() { recoverCallback(w, r, result) },
							r,
						)
					}
					writer.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func doRecoverCallback(doRecoverCallbackFn func(), r *http.Request) {
	defer func() {
		if recoverCallbackRecover := recover(); recoverCallbackRecover != nil {
			internal.Error(r, fmt.Sprintf("%v\n%s", recoverCallbackRecover, debug.Stack()))
		}
	}()

	doRecoverCallbackFn()
}
