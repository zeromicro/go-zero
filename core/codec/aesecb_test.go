package codec

import (
	"crypto/aes"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesEcb(t *testing.T) {
	var (
		key     = []byte("q4t7w!z%C*F-JaNdRgUjXn2r5u8x/A?D")
		val     = []byte("helloworld")
		valLong = []byte("helloworldlong..")
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
	block, err := aes.NewCipher(key)
	assert.NoError(t, err)
	encrypter := NewECBEncrypter(block)
	assert.Equal(t, 16, encrypter.BlockSize())
	decrypter := NewECBDecrypter(block)
	assert.Equal(t, 16, decrypter.BlockSize())

	dst = make([]byte, 8)
	encrypter.CryptBlocks(dst, val)
	for _, b := range dst {
		assert.Equal(t, byte(0), b)
	}

	dst = make([]byte, 8)
	encrypter.CryptBlocks(dst, valLong)
	for _, b := range dst {
		assert.Equal(t, byte(0), b)
	}

	dst = make([]byte, 8)
	decrypter.CryptBlocks(dst, val)
	for _, b := range dst {
		assert.Equal(t, byte(0), b)
	}

	dst = make([]byte, 8)
	decrypter.CryptBlocks(dst, valLong)
	for _, b := range dst {
		assert.Equal(t, byte(0), b)
	}

	_, err = EcbEncryptBase64("cTR0N3dDKkYtSmFOZFJnVWpYbjJyNXU4eC9BP0QK", "aGVsbG93b3JsZGxvbmcuLgo=")
	assert.Error(t, err)
}

func TestAesEcbBase64(t *testing.T) {
	const (
		val     = "hello"
		badKey1 = "aaaaaaaaa"
		// more than 32 chars
		badKey2 = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	)
	key := []byte("q4t7w!z%C*F-JaNdRgUjXn2r5u8x/A?D")
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
	assert.Nil(t, err)
	b, err := base64.StdEncoding.DecodeString(src)
	assert.Nil(t, err)
	assert.Equal(t, val, string(b))
}
