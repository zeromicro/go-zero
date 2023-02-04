package sqlx

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

var slave1 = slave{}

func Test_randomPicker(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		r := newRandomPicker(func() []slave {
			return []slave{slave1}
		})

		s, err := r.pick()
		assert.NoError(t, err)
		assert.NotNil(t, s)
	})

	t.Run("error", func(t *testing.T) {
		r := newRandomPicker(func() []slave {
			return []slave{}
		})
		s, err := r.pick()
		assert.ErrorIs(t, err, errNoAvailableSlave)
		assert.NotNil(t, s)

		r = &randomPicker{
			r:        rand.New(rand.NewSource(1)),
			fnSlaves: nil,
		}
		s, err = r.pick()
		assert.ErrorIs(t, err, errNoAvailableSlave)
		assert.NotNil(t, s)
	})
}

func Test_weightRandomPicker(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		w := newWeightRandomPicker([]int{1, 2, 4}, func() []slave {
			return []slave{
				{
					datasource: "1",
				},
				{
					datasource: "2",
				},
			}
		})
		w.r = rand.New(rand.NewSource(1))

		s, err := w.pick()
		assert.NoError(t, err)
		assert.EqualValues(t, s.datasource, "2")

		s, err = w.pick()
		assert.NoError(t, err)
		assert.EqualValues(t, s.datasource, "2")

		s, err = w.pick()
		assert.NoError(t, err)
		assert.EqualValues(t, s.datasource, "2")
	})

	t.Run("err", func(t *testing.T) {
		w := newWeightRandomPicker([]int{1, 2}, func() []slave {
			return []slave{}
		})
		_, err := w.pick()
		assert.ErrorIs(t, err, errNoAvailableSlave)

		w = newWeightRandomPicker([]int{1, 2}, nil)
		_, err = w.pick()
		assert.ErrorIs(t, err, errNoAvailableSlave)
	})
}

func Test_roundRobinPicker(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		r := newRoundRobinPicker(func() []slave {
			return []slave{
				{
					datasource: "1",
				},
				{
					datasource: "2",
				}, {
					datasource: "3",
				},
			}
		})

		s, err := r.pick()
		assert.NoError(t, err)
		assert.EqualValues(t, s.datasource, "1")

		s, err = r.pick()
		assert.NoError(t, err)
		assert.EqualValues(t, s.datasource, "2")

		s, err = r.pick()
		assert.NoError(t, err)
		assert.EqualValues(t, s.datasource, "3")

		s, err = r.pick()
		assert.NoError(t, err)
		assert.EqualValues(t, s.datasource, "1")

	})

	t.Run("error", func(t *testing.T) {
		r := newRoundRobinPicker(func() []slave {
			return []slave{}
		})
		_, err := r.pick()
		assert.ErrorIs(t, err, errNoAvailableSlave)

		r = newRoundRobinPicker(nil)
		_, err = r.pick()
		assert.ErrorIs(t, err, errNoAvailableSlave)
	})
}

func Test_weightRoundRobinPicker(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		w := newWeightRoundRobinPicker([]int{1, 2, 3, 1}, func() []slave {
			return []slave{
				{
					datasource: "1",
				},
				{
					datasource: "2",
				},
				{
					datasource: "3",
				},
			}
		})
		s, err := w.pick()
		assert.NoError(t, err)
		assert.EqualValues(t, s.datasource, "1")

		for i := 0; i < 2; i++ {
			s, err := w.pick()
			assert.NoError(t, err)
			assert.EqualValues(t, s.datasource, "2")
		}

		for i := 0; i < 3; i++ {
			s, err := w.pick()
			assert.NoError(t, err)
			assert.EqualValues(t, s.datasource, "3")
		}

		s, err = w.pick()
		assert.NoError(t, err)
		assert.EqualValues(t, s.datasource, "3")

		s, err = w.pick()
		assert.NoError(t, err)
		assert.EqualValues(t, s.datasource, "1")

	})
}
