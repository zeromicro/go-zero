package rq

import (
	"math/rand"

	"zero/core/hash"
)

type HashChange struct {
	id      int64
	oldHash *hash.ConsistentHash
	newHash *hash.ConsistentHash
}

func NewHashChange(oldHash, newHash *hash.ConsistentHash) HashChange {
	return HashChange{
		id:      rand.Int63(),
		oldHash: oldHash,
		newHash: newHash,
	}
}

func (hc HashChange) GetId() int64 {
	return hc.id
}

func (hc HashChange) ShallEvict(key interface{}) bool {
	oldTarget, oldOk := hc.oldHash.Get(key)
	if !oldOk {
		return false
	}

	newTarget, newOk := hc.newHash.Get(key)
	if !newOk {
		return false
	}

	return oldTarget != newTarget
}
