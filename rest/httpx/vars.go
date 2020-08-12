package httpx

import "errors"

const (
	ApplicationJson = "application/json"
	ContentEncoding = "Content-Encoding"
	ContentSecurity = "X-Content-Security"
	ContentType     = "Content-Type"
	KeyField        = "key"
	SecretField     = "secret"
	TypeField       = "type"
	CryptionType    = 1
)

const (
	CodeSignaturePass = iota
	CodeSignatureInvalidHeader
	CodeSignatureWrongTime
	CodeSignatureInvalidToken
)

var ErrBodylessRequest = errors.New("not a POST|PUT|PATCH request")
