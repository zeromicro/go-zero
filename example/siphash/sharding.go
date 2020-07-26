package sharding

import "github.com/dchest/siphash"

func sharding(token string) uint64 {
	sum := siphash.Hash(0, 0, []byte(token))
	return sum % 3
}
