package codec

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// ErrPaddingSize indicates bad padding size.
var ErrPaddingSize = errors.New("padding size error")

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// Deprecated: NewECBEncrypter returns an ECB encrypter.
// ECB mode is insecure for multi-block data. Use AES-GCM instead.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

// BlockSize returns the mode's block size.
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }

// CryptBlocks encrypts a number of blocks. The length of src must be a multiple of
// the block size. Dst and src must overlap entirely or not at all.
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}

	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// Deprecated: NewECBDecrypter returns an ECB decrypter.
// ECB mode is insecure for multi-block data. Use AES-GCM instead.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

// BlockSize returns the mode's block size.
func (x *ecbDecrypter) BlockSize() int {
	return x.blockSize
}

// CryptBlocks decrypts a number of blocks. The length of src must be a multiple of
// the block size. Dst and src must overlap entirely or not at all.
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}

	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

// Deprecated: EcbDecrypt decrypts src with the given key.
// ECB mode is insecure for multi-block data. Use AES-GCM instead.
func EcbDecrypt(key, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(src)%block.BlockSize() != 0 {
		return nil, ErrPaddingSize
	}

	decrypter := NewECBDecrypter(block)
	decrypted := make([]byte, len(src))
	decrypter.CryptBlocks(decrypted, src)

	return pkcs5Unpadding(decrypted, decrypter.BlockSize())
}

// Deprecated: EcbDecryptBase64 decrypts base64 encoded src with the given base64 encoded key.
// The returned string is also base64 encoded.
// ECB mode is insecure for multi-block data. Use AES-GCM instead.
func EcbDecryptBase64(key, src string) (string, error) {
	keyBytes, err := getKeyBytes(key)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	decryptedBytes, err := EcbDecrypt(keyBytes, encryptedBytes)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(decryptedBytes), nil
}

// Deprecated: EcbEncrypt encrypts src with the given key.
// ECB mode is insecure for multi-block data. Use AES-GCM instead.
func EcbEncrypt(key, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	padded := pkcs5Padding(src, block.BlockSize())
	crypted := make([]byte, len(padded))
	encrypter := NewECBEncrypter(block)
	encrypter.CryptBlocks(crypted, padded)

	return crypted, nil
}

// Deprecated: EcbEncryptBase64 encrypts base64 encoded src with the given base64 encoded key.
// The returned string is also base64 encoded.
// ECB mode is insecure for multi-block data. Use AES-GCM instead.
func EcbEncryptBase64(key, src string) (string, error) {
	keyBytes, err := getKeyBytes(key)
	if err != nil {
		return "", err
	}

	srcBytes, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := EcbEncrypt(keyBytes, srcBytes)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

func getKeyBytes(key string) ([]byte, error) {
	if len(key) <= 32 {
		return []byte(key), nil
	}

	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	return keyBytes, nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5Unpadding(src []byte, blockSize int) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, ErrPaddingSize
	}

	unpadding := int(src[length-1])
	if unpadding < 1 || unpadding > blockSize || unpadding > length {
		return nil, ErrPaddingSize
	}

	for _, b := range src[length-unpadding:] {
		if int(b) != unpadding {
			return nil, ErrPaddingSize
		}
	}

	return src[:length-unpadding], nil
}
