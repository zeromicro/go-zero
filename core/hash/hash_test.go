package hash

import (
	"crypto/md5"
	"fmt"
	"hash/fnv"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	text      = "hello, world!\n"
	md5Digest = "910c8bc73110b0cd1bc5d2bcae782511"
)

func TestMd5(t *testing.T) {
	actual := fmt.Sprintf("%x", Md5([]byte(text)))
	assert.Equal(t, md5Digest, actual)
}

func BenchmarkHashFnv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := fnv.New32()
		new(big.Int).SetBytes(h.Sum([]byte(text))).Int64()
	}
}

func BenchmarkHashMd5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := md5.New()
		bytes := h.Sum([]byte(text))
		new(big.Int).SetBytes(bytes).Int64()
	}
}

func BenchmarkMurmur3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Hash([]byte(text))
	}
}
