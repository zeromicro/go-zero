package handler

import (
	"net/http"
	"time"

	"github.com/tal-tech/go-zero/core/codec"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest/httpx"
	"github.com/tal-tech/go-zero/rest/internal/security"
)

const contentSecurity = "X-Content-Security"

type UnsignedCallback func(w http.ResponseWriter, r *http.Request, next http.Handler, strict bool, code int)

func ContentSecurityHandler(decrypters map[string]codec.RsaDecrypter, tolerance time.Duration,
	strict bool, callbacks ...UnsignedCallback) func(http.Handler) http.Handler {
	if len(callbacks) == 0 {
		callbacks = append(callbacks, handleVerificationFailure)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodDelete, http.MethodGet, http.MethodPost, http.MethodPut:
				header, err := security.ParseContentSecurity(decrypters, r)
				if err != nil {
					logx.Infof("Signature parse failed, X-Content-Security: %s, error: %s",
						r.Header.Get(contentSecurity), err.Error())
					executeCallbacks(w, r, next, strict, httpx.CodeSignatureInvalidHeader, callbacks)
				} else if code := security.VerifySignature(r, header, tolerance); code != httpx.CodeSignaturePass {
					logx.Infof("Signature verification failed, X-Content-Security: %s",
						r.Header.Get(contentSecurity))
					executeCallbacks(w, r, next, strict, code, callbacks)
				} else if r.ContentLength > 0 && header.Encrypted() {
					CryptionHandler(header.Key)(next).ServeHTTP(w, r)
				} else {
					next.ServeHTTP(w, r)
				}
			default:
				next.ServeHTTP(w, r)
			}
		})
	}
}

func executeCallbacks(w http.ResponseWriter, r *http.Request, next http.Handler, strict bool,
	code int, callbacks []UnsignedCallback) {
	for _, callback := range callbacks {
		callback(w, r, next, strict, code)
	}
}

func handleVerificationFailure(w http.ResponseWriter, r *http.Request, next http.Handler, strict bool, code int) {
	if strict {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		next.ServeHTTP(w, r)
	}
}
