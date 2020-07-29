package update

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"sort"

	"zero/core/hash"
	"zero/core/jsonx"
	"zero/rq/internal"
)

var ErrInvalidServerChange = errors.New("not a server change message")

type (
	weightedKey struct {
		Key    string
		Weight int
	}

	Snapshot struct {
		Keys         []string
		WeightedKeys []weightedKey
	}

	ServerChange struct {
		Previous Snapshot
		Current  Snapshot
		Servers  []string
	}
)

func (s Snapshot) GetCode() string {
	keys := append([]string(nil), s.Keys...)
	sort.Strings(keys)
	weightedKeys := append([]weightedKey(nil), s.WeightedKeys...)
	sort.SliceStable(weightedKeys, func(i, j int) bool {
		return weightedKeys[i].Key < weightedKeys[j].Key
	})

	digest := md5.New()
	for _, key := range keys {
		io.WriteString(digest, fmt.Sprintf("%s\n", key))
	}
	for _, wkey := range weightedKeys {
		io.WriteString(digest, fmt.Sprintf("%s:%d\n", wkey.Key, wkey.Weight))
	}

	return fmt.Sprintf("%x", digest.Sum(nil))
}

func (sc ServerChange) CreateCurrentHash() *hash.ConsistentHash {
	curHash := hash.NewConsistentHash()

	for _, key := range sc.Current.Keys {
		curHash.Add(key)
	}
	for _, wkey := range sc.Current.WeightedKeys {
		curHash.AddWithWeight(wkey.Key, wkey.Weight)
	}

	return curHash
}

func (sc ServerChange) CreatePrevHash() *hash.ConsistentHash {
	prevHash := hash.NewConsistentHash()

	for _, key := range sc.Previous.Keys {
		prevHash.Add(key)
	}
	for _, wkey := range sc.Previous.WeightedKeys {
		prevHash.AddWithWeight(wkey.Key, wkey.Weight)
	}

	return prevHash
}

func (sc ServerChange) GetCode() string {
	return sc.Current.GetCode()
}

func IsServerChange(message string) bool {
	return len(message) > 0 && message[0] == internal.ServerSensitivePrefix
}

func (sc ServerChange) Marshal() (string, error) {
	body, err := jsonx.Marshal(sc)
	if err != nil {
		return "", err
	}

	return string(append([]byte{internal.ServerSensitivePrefix}, body...)), nil
}

func UnmarshalServerChange(body string) (ServerChange, error) {
	if len(body) == 0 {
		return ServerChange{}, ErrInvalidServerChange
	}

	var change ServerChange
	err := jsonx.UnmarshalFromString(body[1:], &change)

	return change, err
}
