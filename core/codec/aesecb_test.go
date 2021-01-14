package codec

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesEcb(t *testing.T) {
	var (
		key     = []byte("q4t7w!z%C*F-JaNdRgUjXn2r5u8x/A?D")
		val     = []byte("hello")
		badKey1 = []byte("aaaaaaaaa")
		// more than 32 chars
		badKey2 = []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	)
	_, err := EcbEncrypt(badKey1, val)
	assert.NotNil(t, err)
	_, err = EcbEncrypt(badKey2, val)
	assert.NotNil(t, err)
	dst, err := EcbEncrypt(key, val)
	assert.Nil(t, err)
	_, err = EcbDecrypt(badKey1, dst)
	assert.NotNil(t, err)
	_, err = EcbDecrypt(badKey2, dst)
	assert.NotNil(t, err)
	_, err = EcbDecrypt(key, val)
	// not enough block, just nil
	assert.Nil(t, err)
	src, err := EcbDecrypt(key, dst)
	assert.Nil(t, err)
	assert.Equal(t, val, src)
}

func TestAesEcbBase64(t *testing.T) {
	const (
		val     = "hello"
		badKey1 = "aaaaaaaaa"
		// more than 32 chars
		badKey2 = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	)
	var key = []byte("q4t7w!z%C*F-JaNdRgUjXn2r5u8x/A?D")
	b64Key := base64.StdEncoding.EncodeToString(key)
	b64Val := base64.StdEncoding.EncodeToString([]byte(val))
	_, err := EcbEncryptBase64(badKey1, val)
	assert.NotNil(t, err)
	_, err = EcbEncryptBase64(badKey2, val)
	assert.NotNil(t, err)
	_, err = EcbEncryptBase64(b64Key, val)
	assert.NotNil(t, err)
	dst, err := EcbEncryptBase64(b64Key, b64Val)
	assert.Nil(t, err)
	_, err = EcbDecryptBase64(badKey1, dst)
	assert.NotNil(t, err)
	_, err = EcbDecryptBase64(badKey2, dst)
	assert.NotNil(t, err)
	_, err = EcbDecryptBase64(b64Key, val)
	assert.NotNil(t, err)
	src, err := EcbDecryptBase64(b64Key, dst)
	b, err := base64.StdEncoding.DecodeString(src)
	assert.Nil(t, err)
	assert.Equal(t, val, string(b))
}
