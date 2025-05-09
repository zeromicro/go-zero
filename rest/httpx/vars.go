package httpx

import "github.com/zeromicro/go-zero/rest/internal/header"

const (
	// ContentEncoding means Content-Encoding.
	ContentEncoding = "Content-Encoding"
	// ContentSecurity means X-Content-Security.
	ContentSecurity = "X-Content-Security"
	// ContentType means Content-Type.
	ContentType = header.ContentType
	// JsonContentType means application/json.
	JsonContentType = header.ContentTypeJson
	// KeyField means key.
	KeyField = "key"
	// SecretField means secret.
	SecretField = "secret"
	// TypeField means type.
	TypeField = "type"
	// CryptionType means cryption.
	CryptionType = 1
)

const (
	// CodeSignaturePass means signature verification passed.
	CodeSignaturePass = iota
	// CodeSignatureInvalidHeader means invalid header in signature.
	CodeSignatureInvalidHeader
	// CodeSignatureWrongTime means wrong timestamp in signature.
	CodeSignatureWrongTime
	// CodeSignatureInvalidToken means invalid token in signature.
	CodeSignatureInvalidToken
)
