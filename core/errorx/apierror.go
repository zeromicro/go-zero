package errorx

import (
	"net/http"
)

type SimpleMsg struct {
	Msg string `json:"msg"`
}

// ApiError is error with http status codes.
type ApiError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *ApiError) Error() string {
	return e.Msg
}

// NewApiError returns Api Error
func NewApiError(code int, msg string) error {
	return &ApiError{Code: code, Msg: msg}
}

// NewApiErrorWithoutMsg returns Api Error without message
func NewApiErrorWithoutMsg(code int) error {
	return &ApiError{Code: code, Msg: ""}
}

// NewApiInternalError returns Api Error with http internal error status code
func NewApiInternalError(msg string) error {
	return &ApiError{Code: http.StatusInternalServerError, Msg: msg}
}

// NewApiBadRequestError returns Api Error with http bad request status code
func NewApiBadRequestError(msg string) error {
	return &ApiError{Code: http.StatusBadRequest, Msg: msg}
}

// NewApiUnauthorizedError returns Api Error with http unauthorized status code
func NewApiUnauthorizedError(msg string) error {
	return &ApiError{Code: http.StatusUnauthorized, Msg: msg}
}

// NewApiForbiddenError returns Api Error with http forbidden status code
func NewApiForbiddenError(msg string) error {
	return &ApiError{Code: http.StatusForbidden, Msg: msg}
}

// NewApiNotFoundError returns Api Error with http not found status code
func NewApiNotFoundError(msg string) error {
	return &ApiError{Code: http.StatusNotFound, Msg: msg}
}

// NewApiBadGatewayError returns Api Error with http bad gateway status code
func NewApiBadGatewayError(msg string) error {
	return &ApiError{Code: http.StatusBadGateway, Msg: msg}
}
