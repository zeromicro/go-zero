package sqlx

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

var (
	errNoAvailableSlave = errors.New("no available slave")

	emptySlave = slave{}
)

type (
	picker interface {
		pick() (slave, error)
	}
	fnSlaves func() []slave
)

type randomPicker struct {
	r        *rand.Rand
	fnSlaves fnSlaves
}

func newRandomPicker(fn fnSlaves) *randomPicker {
	return &randomPicker{r: rand.New(rand.NewSource(time.Now().UnixNano())), fnSlaves: fn}
}

func (r *randomPicker) pick() (slave, error) {
	if r.fnSlaves == nil {
		return emptySlave, errNoAvailableSlave
	}

	slaves := r.fnSlaves()
	if len(slaves) == 0 {
		return emptySlave, errNoAvailableSlave
	}

	return slaves[r.r.Intn(len(slaves))], nil
}

type weightRandomPicker struct {
	r        *rand.Rand
	weights  []int
	fnSlaves func() []slave
}

func newWeightRandomPicker(weights []int, fn fnSlaves) *weightRandomPicker {
	return &weightRandomPicker{weights: weights, r: rand.New(rand.NewSource(time.Now().UnixNano())), fnSlaves: fn}
}

func (w *weightRandomPicker) pick() (slave, error) {
	if w.fnSlaves == nil {
		return emptySlave, errNoAvailableSlave
	}

	slaves := w.fnSlaves()
	if len(slaves) == 0 {
		return emptySlave, errNoAvailableSlave
	}

	var weightRands = make([]int, 0, len(w.weights))
	for i := 0; i < len(w.weights); i++ {
		for n := 0; n < w.weights[i]; n++ {
			weightRands = append(weightRands, i)
		}
	}

	index := weightRands[w.r.Intn(len(weightRands))]
	if index >= len(slaves) {
		index = len(slaves) - 1
	}

	return slaves[index], nil
}

type roundRobinPicker struct {
	i        int
	fnSlaves func() []slave
	mu       sync.Mutex
}

func newRoundRobinPicker(fn fnSlaves) *roundRobinPicker {
	return &roundRobinPicker{
		fnSlaves: fn,
	}
}

func (r *roundRobinPicker) pick() (slave, error) {
	if r.fnSlaves == nil {
		return emptySlave, errNoAvailableSlave
	}

	slaves := r.fnSlaves()
	if len(slaves) == 0 {
		return emptySlave, errNoAvailableSlave
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.i >= len(slaves) {
		r.i = 0
	}

	return slaves[r.i], nil
}
