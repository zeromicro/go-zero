package resolver

import (
	"github.com/tal-tech/go-zero/core/hash"
	"github.com/tal-tech/go-zero/core/sysx"
	"math/rand"
)

func subset(set []string, sub int) []string {
	if len(set) <= sub {
		rand.Shuffle(len(set), func(i, j int) {
			set[i], set[j] = set[j], set[i]
		})
		return set
	}

	// group clients into rounds, each round uses the same shuffled list
	count := uint64(len(set) / sub)
	clientID := hash.Hash([]byte(sysx.Hostname()))
	round := clientID / count

	r := rand.New(rand.NewSource(int64(round)))
	r.Shuffle(len(set), func(i, j int) {
		set[i], set[j] = set[j], set[i]
	})

	start := (clientID % count) * uint64(sub)
	return set[start : start+uint64(sub)]
}
