package errorx

import (
	"net/http"
)

type SimpleMsg struct {
	Msg string `json:"msg"`
}

type ApiError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *ApiError) Error() string {
	return e.Msg
}

func NewApiError(code int, msg string) error {
	return &ApiError{Code: code, Msg: msg}
}

func NewApiErrorWithoutMsg(code int) error {
	return &ApiError{Code: code, Msg: ""}
}

func NewApiInternalError(msg string) error {
	return &ApiError{Code: http.StatusInternalServerError, Msg: msg}
}

func NewApiBadRequestError(msg string) error {
	return &ApiError{Code: http.StatusBadRequest, Msg: msg}
}

func NewApiUnauthorizedError(msg string) error {
	return &ApiError{Code: http.StatusUnauthorized, Msg: msg}
}

func NewApiForbiddenError(msg string) error {
	return &ApiError{Code: http.StatusForbidden, Msg: msg}
}

func NewApiNotFoundError(msg string) error {
	return &ApiError{Code: http.StatusNotFound, Msg: msg}
}

func NewApiBadGatewayError(msg string) error {
	return &ApiError{Code: http.StatusBadGateway, Msg: msg}
}
