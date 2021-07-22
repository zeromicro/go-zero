package hash

import (
	"crypto/md5"
	"fmt"

	"github.com/spaolacci/murmur3"
)

// Hash returns the hash value of data.
func Hash(data []byte) uint64 {
	return murmur3.Sum64(data)
}

// Md5 returns the md5 bytes of data.
func Md5(data []byte) []byte {
	digest := md5.New()
	digest.Write(data)
	return digest.Sum(nil)
}

// Md5Hex returns the md5 hex string of data.
func Md5Hex(data []byte) string {
	return fmt.Sprintf("%x", Md5(data))
}
