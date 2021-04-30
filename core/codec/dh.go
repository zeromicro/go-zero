package codec

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// see https://www.zhihu.com/question/29383090/answer/70435297
// see https://www.ietf.org/rfc/rfc3526.txt
// 2048-bit MODP Group

var (
	// ErrInvalidPriKey indicates the invalid private key.
	ErrInvalidPriKey = errors.New("invalid private key")
	// ErrInvalidPubKey indicates the invalid public key.
	ErrInvalidPubKey = errors.New("invalid public key")
	// ErrPubKeyOutOfBound indicates the public key is out of bound.
	ErrPubKeyOutOfBound = errors.New("public key out of bound")

	p, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF6955817183995497CEA956AE515D2261898FA051015728E5A8AACAA68FFFFFFFFFFFFFFFF", 16)
	g, _ = new(big.Int).SetString("2", 16)
	zero = big.NewInt(0)
)

// DhKey defines the Diffie Hellman key.
type DhKey struct {
	PriKey *big.Int
	PubKey *big.Int
}

// ComputeKey returns a key from public key and private key.
func ComputeKey(pubKey, priKey *big.Int) (*big.Int, error) {
	if pubKey == nil {
		return nil, ErrInvalidPubKey
	}

	if pubKey.Sign() <= 0 && p.Cmp(pubKey) <= 0 {
		return nil, ErrPubKeyOutOfBound
	}

	if priKey == nil {
		return nil, ErrInvalidPriKey
	}

	return new(big.Int).Exp(pubKey, priKey, p), nil
}

// GenerateKey returns a Diffie Hellman key.
func GenerateKey() (*DhKey, error) {
	var err error
	var x *big.Int

	for {
		x, err = rand.Int(rand.Reader, p)
		if err != nil {
			return nil, err
		}

		if zero.Cmp(x) < 0 {
			break
		}
	}

	key := new(DhKey)
	key.PriKey = x
	key.PubKey = new(big.Int).Exp(g, x, p)

	return key, nil
}

// NewPublicKey returns a public key from the given bytes.
func NewPublicKey(bs []byte) *big.Int {
	return new(big.Int).SetBytes(bs)
}

// Bytes returns public key bytes.
func (k *DhKey) Bytes() []byte {
	if k.PubKey == nil {
		return nil
	}

	byteLen := (p.BitLen() + 7) >> 3
	ret := make([]byte, byteLen)
	copyWithLeftPad(ret, k.PubKey.Bytes())

	return ret
}

func copyWithLeftPad(dst, src []byte) {
	padBytes := len(dst) - len(src)
	for i := 0; i < padBytes; i++ {
		dst[i] = 0
	}
	copy(dst[padBytes:], src)
}
