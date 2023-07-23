package codec

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffieHellman(t *testing.T) {
	key1, err := GenerateKey()
	assert.Nil(t, err)
	key2, err := GenerateKey()
	assert.Nil(t, err)

	pubKey1, err := ComputeKey(key1.PubKey, key2.PriKey)
	assert.Nil(t, err)
	pubKey2, err := ComputeKey(key2.PubKey, key1.PriKey)
	assert.Nil(t, err)

	assert.Equal(t, pubKey1, pubKey2)
}

func TestDiffieHellman1024(t *testing.T) {
	old := p
	p, _ = new(big.Int).SetString("F488FD584E49DBCD20B49DE49107366B336C380D451D0F7C88B31C7C5B2D8EF6F3C923C043F0A55B188D8EBB558CB85D38D334FD7C175743A31D186CDE33212CB52AFF3CE1B1294018118D7C84A70A72D686C40319C807297ACA950CD9969FABD00A509B0246D3083D66A45D419F9C7CBD894B221926BAABA25EC355E92F78C7", 16)
	defer func() {
		p = old
	}()

	key1, err := GenerateKey()
	assert.Nil(t, err)
	key2, err := GenerateKey()
	assert.Nil(t, err)

	pubKey1, err := ComputeKey(key1.PubKey, key2.PriKey)
	assert.Nil(t, err)
	pubKey2, err := ComputeKey(key2.PubKey, key1.PriKey)
	assert.Nil(t, err)

	assert.Equal(t, pubKey1, pubKey2)
}

func TestDiffieHellmanMiddleManAttack(t *testing.T) {
	key1, err := GenerateKey()
	assert.Nil(t, err)
	keyMiddle, err := GenerateKey()
	assert.Nil(t, err)
	key2, err := GenerateKey()
	assert.Nil(t, err)

	const aesByteLen = 32
	pubKey1, err := ComputeKey(keyMiddle.PubKey, key1.PriKey)
	assert.Nil(t, err)
	src := []byte(`hello, world!`)
	encryptedSrc, err := EcbEncrypt(pubKey1.Bytes()[:aesByteLen], src)
	assert.Nil(t, err)
	pubKeyMiddle, err := ComputeKey(key1.PubKey, keyMiddle.PriKey)
	assert.Nil(t, err)
	decryptedSrc, err := EcbDecrypt(pubKeyMiddle.Bytes()[:aesByteLen], encryptedSrc)
	assert.Nil(t, err)
	assert.Equal(t, string(src), string(decryptedSrc))

	pubKeyMiddle, err = ComputeKey(key2.PubKey, keyMiddle.PriKey)
	assert.Nil(t, err)
	encryptedSrc, err = EcbEncrypt(pubKeyMiddle.Bytes()[:aesByteLen], decryptedSrc)
	assert.Nil(t, err)
	pubKey2, err := ComputeKey(keyMiddle.PubKey, key2.PriKey)
	assert.Nil(t, err)
	decryptedSrc, err = EcbDecrypt(pubKey2.Bytes()[:aesByteLen], encryptedSrc)
	assert.Nil(t, err)
	assert.Equal(t, string(src), string(decryptedSrc))
}

func TestKeyBytes(t *testing.T) {
	var empty DhKey
	assert.Equal(t, 0, len(empty.Bytes()))

	key, err := GenerateKey()
	assert.Nil(t, err)
	assert.True(t, len(key.Bytes()) > 0)
}

func TestDHOnErrors(t *testing.T) {
	key, err := GenerateKey()
	assert.Nil(t, err)
	assert.NotEmpty(t, key.Bytes())
	_, err = ComputeKey(key.PubKey, key.PriKey)
	assert.NoError(t, err)
	_, err = ComputeKey(nil, key.PriKey)
	assert.Error(t, err)
	_, err = ComputeKey(key.PubKey, nil)
	assert.Error(t, err)

	assert.NotNil(t, NewPublicKey([]byte("")))
}
