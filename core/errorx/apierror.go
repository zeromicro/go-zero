package errorx

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
