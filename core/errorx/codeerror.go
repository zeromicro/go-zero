package errorx

const defaultCode = 3

// CodeError is error with custom error code
type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type CodeErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// NewCodeError returns a code error
func NewCodeError(code int, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

// NewCodeCanceledError returns Code Error with custom cancel error code
func NewCodeCanceledError(msg string) error {
	return &CodeError{Code: 1, Msg: msg}
}

// NewCodeInvalidArgumentError returns Code Error with custom invalid argument error code
func NewCodeInvalidArgumentError(msg string) error {
	return &CodeError{Code: 3, Msg: msg}
}

// NewCodeNotFoundError returns Code Error with custom not found error code
func NewCodeNotFoundError(msg string) error {
	return &CodeError{Code: 5, Msg: msg}
}

// NewCodeAlreadyExistsError returns Code Error with custom already exists error code
func NewCodeAlreadyExistsError(msg string) error {
	return &CodeError{Code: 6, Msg: msg}
}

// NewCodeAbortedError returns Code Error with custom aborted error code
func NewCodeAbortedError(msg string) error {
	return &CodeError{Code: 10, Msg: msg}
}

// NewCodeInternalError returns Code Error with custom internal error code
func NewCodeInternalError(msg string) error {
	return &CodeError{Code: 13, Msg: msg}
}

// NewCodeUnavailableError returns Code Error with custom unavailable error code
func NewCodeUnavailableError(msg string) error {
	return &CodeError{Code: 14, Msg: msg}
}

func NewDefaultError(msg string) error {
	return NewCodeError(defaultCode, msg)
}

func (e *CodeError) Error() string {
	return e.Msg
}

func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: e.Code,
		Msg:  e.Msg,
	}
}
