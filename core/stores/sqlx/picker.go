package sqlx

import (
	"math/rand"
	"time"
)

type picker interface {
	pick() slave
}

type randomPicker struct {
	r      *rand.Rand
	slaves []slave
}

func newRandomPicker(slaves []slave) *randomPicker {
	return &randomPicker{r: rand.New(rand.NewSource(time.Now().UnixNano())), slaves: slaves}
}

func (r *randomPicker) pick() slave {
	return r.slaves[r.r.Intn(len(r.slaves))]
}
