package bloom

import (
	"errors"
	"strconv"

	"github.com/tal-tech/go-zero/core/hash"
	"github.com/tal-tech/go-zero/core/stores/redis"
)

const (
	// for detailed error rate table, see http://pages.cs.wisc.edu/~cao/papers/summary-cache/node8.html
	// maps as k in the error rate table
	maps      = 14
	setScript = `
local key = KEYS[1]
for _, offset in ipairs(ARGV) do
	redis.call("setbit", key, offset, 1)
end
`
	testScript = `
local key = KEYS[1]
for _, offset in ipairs(ARGV) do
	if tonumber(redis.call("getbit", key, offset)) == 0 then
		return false
	end
end
return true
`
)

var ErrTooLargeOffset = errors.New("too large offset")

type (
	BitSetProvider interface {
		check([]uint) (bool, error)
		set([]uint) error
	}

	BloomFilter struct {
		bits   uint
		maps   uint
		bitSet BitSetProvider
	}
)

// New create a BloomFilter, store is the backed redis, key is the key for the bloom filter,
// bits is how many bits will be used, maps is how many hashes for each addition.
// best practices:
// elements - means how many actual elements
// when maps = 14, formula: 0.7*(bits/maps), bits = 20*elements, the error rate is 0.000067 < 1e-4
// for detailed error rate table, see http://pages.cs.wisc.edu/~cao/papers/summary-cache/node8.html
func New(store *redis.Redis, key string, bits uint) *BloomFilter {
	return &BloomFilter{
		bits:   bits,
		bitSet: newRedisBitSet(store, key, bits),
	}
}

func (f *BloomFilter) Add(data []byte) error {
	locations := f.getLocations(data)
	err := f.bitSet.set(locations)
	if err != nil {
		return err
	}
	return nil
}

func (f *BloomFilter) Exists(data []byte) (bool, error) {
	locations := f.getLocations(data)
	isSet, err := f.bitSet.check(locations)
	if err != nil {
		return false, err
	}
	if !isSet {
		return false, nil
	}

	return true, nil
}

func (f *BloomFilter) getLocations(data []byte) []uint {
	locations := make([]uint, maps)
	for i := uint(0); i < maps; i++ {
		hashValue := hash.Hash(append(data, byte(i)))
		locations[i] = uint(hashValue % uint64(f.bits))
	}

	return locations
}

type redisBitSet struct {
	store *redis.Redis
	key   string
	bits  uint
}

func newRedisBitSet(store *redis.Redis, key string, bits uint) *redisBitSet {
	return &redisBitSet{
		store: store,
		key:   key,
		bits:  bits,
	}
}

func (r *redisBitSet) buildOffsetArgs(offsets []uint) ([]string, error) {
	var args []string

	for _, offset := range offsets {
		if offset >= r.bits {
			return nil, ErrTooLargeOffset
		}

		args = append(args, strconv.FormatUint(uint64(offset), 10))
	}

	return args, nil
}

func (r *redisBitSet) check(offsets []uint) (bool, error) {
	args, err := r.buildOffsetArgs(offsets)
	if err != nil {
		return false, err
	}

	resp, err := r.store.Eval(testScript, []string{r.key}, args)
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if exists, ok := resp.(int64); !ok {
		return false, nil
	} else {
		return exists == 1, nil
	}
}

func (r *redisBitSet) del() error {
	_, err := r.store.Del(r.key)
	return err
}

func (r *redisBitSet) expire(seconds int) error {
	return r.store.Expire(r.key, seconds)
}

func (r *redisBitSet) set(offsets []uint) error {
	args, err := r.buildOffsetArgs(offsets)
	if err != nil {
		return err
	}

	_, err = r.store.Eval(setScript, []string{r.key}, args)
	if err == redis.Nil {
		return nil
	} else {
		return err
	}
}
