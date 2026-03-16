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

func TestMd5Hex(t *testing.T) {
	actual := Md5Hex([]byte(text))
	assert.Equal(t, md5Digest, actual)
}

func TestHash(t *testing.T) {
	result := Hash([]byte(text))
	assert.NotEqual(t, uint64(0), result)
}

func TestHash_Deterministic(t *testing.T) {
	data := []byte("consistent-hash-test")
	first := Hash(data)
	second := Hash(data)
	assert.Equal(t, first, second)
}

func TestHash_Empty(t *testing.T) {
	// Hash should not panic on empty input.
	result := Hash([]byte{})
	_ = result
}

func TestMd5Hex_Empty(t *testing.T) {
	result := Md5Hex([]byte{})
	assert.Equal(t, 32, len(result))
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
