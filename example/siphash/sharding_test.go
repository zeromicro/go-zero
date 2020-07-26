package sharding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSiphash64(t *testing.T) {
	users := [][]string{
		{
			"5a4b7347200a6e0c185d6101",
			"5b74c444acdd315c509b78fe",
			"5c03e009a496130c2d9bc970",
			"5c6ab5a74867f267d560dd9f",
			"5b80a2b28be129507d176284",
			"5a4b7347200a6e0c185d6101",
			"5b74c444acdd315c509b78fe",
			"5c03e009a496130c2d9bc970",
			"5c6ab5a74867f267d560dd9f",
			"5b80a2b28be129507d176284",
			"5b8d157aacdd313508a892f2",
			"5bf942b4a496130c2d9b7378",
			"5c7fc28cd065f17f9edd3698",
			"5bf40bd22c64fc5ea63a5174",
		},
		{
			"5b839929acdd31271f03ded5",
			"5bc9e28e2c64fc1a69a28e36",
			"5b935d96a49613677b90b589",
			"5b97acb2a49613677b910f47",
			"5c902f3aff5be73689b4b522",
		},
		{
			"5cdbee881a722f0001b9ce99",
			"",
			"5caca58f53add40001c20aaa",
			"5beee68520c25041544e353a",
			"5b0b957d0179b05769cbecde",
			"5bbf45940ab7b7589aa1025f",
			"5ac63009200a6e79cadf5175",
			"5c94ed250ab7b7386c294662",
			"5b9f8ccb2c64fc5832e47d3f",
		},
	}

	for shard, ids := range users {
		for _, id := range ids {
			assert.Equal(t, uint64(shard), sharding(id))
		}
	}
}
