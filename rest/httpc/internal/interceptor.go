package internal

import "net/http"

type (
	Interceptor     func(r *http.Request) (*http.Request, ResponseHandler)
	ResponseHandler func(resp *http.Response, err error)
)
