package test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Data[T, Y any] struct {
	Name  string
	Input T
	Want  Y
	E     error
}

type Option[T, Y any] func(*Executor[T, Y])
type assertFn[Y any] func(t *testing.T, expected, actual Y) bool

func WithComparison[T, Y any](comparisonFn assertFn[Y]) Option[T, Y] {
	return func(e *Executor[T, Y]) {
		e.equalFn = comparisonFn
	}
}

type Executor[T, Y any] struct {
	list    []Data[T, Y]
	equalFn assertFn[Y]
}

func NewExecutor[T, Y any](opt ...Option[T, Y]) *Executor[T, Y] {
	e := &Executor[T, Y]{}
	opt = append(opt, WithComparison[T, Y](func(t *testing.T, expected, actual Y) bool {
		gotBytes, err := json.Marshal(actual)
		if err != nil {
			t.Fatal(err)
			return false
		}
		wantBytes, err := json.Marshal(expected)
		if err != nil {
			t.Fatal(err)
			return false
		}
		return assert.JSONEq(t, string(wantBytes), string(gotBytes))
	}))

	for _, o := range opt {
		o(e)
	}
	return e
}

func (e *Executor[T, Y]) Add(data ...Data[T, Y]) {
	e.list = append(e.list, data...)
}

func (e *Executor[T, Y]) Run(t *testing.T, do func(T) Y) {
	if do == nil {
		panic("execution body is nil")
		return
	}
	for _, v := range e.list {
		t.Run(v.Name, func(t *testing.T) {
			inner := do
			e.equalFn(t, v.Want, inner(v.Input))
		})
	}
}

func (e *Executor[T, Y]) RunE(t *testing.T, do func(T) (Y, error)) {
	if do == nil {
		panic("execution body is nil")
		return
	}
	for _, v := range e.list {
		t.Run(v.Name, func(t *testing.T) {
			inner := do
			got, err := inner(v.Input)
			if v.E != nil {
				assert.Equal(t, v.E, err)
				return
			}
			e.equalFn(t, v.Want, got)
		})
	}
}
