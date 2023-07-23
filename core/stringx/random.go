package stringx

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits  = 6 // 6 bits to represent a letter index
	idLen          = 8
	defaultRandLen = 8
	letterIdxMask  = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax   = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = newLockedSource(time.Now().UnixNano())

type lockedSource struct {
	source rand.Source
	lock   sync.Mutex
}

func newLockedSource(seed int64) *lockedSource {
	return &lockedSource{
		source: rand.NewSource(seed),
	}
}

func (ls *lockedSource) Int63() int64 {
	ls.lock.Lock()
	defer ls.lock.Unlock()
	return ls.source.Int63()
}

func (ls *lockedSource) Seed(seed int64) {
	ls.lock.Lock()
	defer ls.lock.Unlock()
	ls.source.Seed(seed)
}

// Rand returns a random string.
func Rand() string {
	return Randn(defaultRandLen)
}

// RandId returns a random id string.
func RandId() string {
	b := make([]byte, idLen)
	_, err := crand.Read(b)
	if err != nil {
		return Randn(idLen)
	}

	return fmt.Sprintf("%x%x%x%x", b[0:2], b[2:4], b[4:6], b[6:8])
}

// Randn returns a random string with length n.
func Randn(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// Seed sets the seed to seed.
func Seed(seed int64) {
	src.Seed(seed)
}
