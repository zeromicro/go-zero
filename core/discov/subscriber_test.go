package discov

import (
	"testing"

	"zero/core/discov/internal"

	"github.com/stretchr/testify/assert"
)

const (
	actionAdd = iota
	actionDel
)

func TestContainer(t *testing.T) {
	type action struct {
		act int
		key string
		val string
	}
	tests := []struct {
		name   string
		do     []action
		expect []string
	}{
		{
			name: "add one",
			do: []action{
				{
					act: actionAdd,
					key: "first",
					val: "a",
				},
			},
			expect: []string{
				"a",
			},
		},
		{
			name: "add two",
			do: []action{
				{
					act: actionAdd,
					key: "first",
					val: "a",
				},
				{
					act: actionAdd,
					key: "second",
					val: "b",
				},
			},
			expect: []string{
				"a",
				"b",
			},
		},
		{
			name: "add two, delete one",
			do: []action{
				{
					act: actionAdd,
					key: "first",
					val: "a",
				},
				{
					act: actionAdd,
					key: "second",
					val: "b",
				},
				{
					act: actionDel,
					key: "first",
					val: "a",
				},
			},
			expect: []string{
				"b",
			},
		},
		{
			name: "add two, delete two",
			do: []action{
				{
					act: actionAdd,
					key: "first",
					val: "a",
				},
				{
					act: actionAdd,
					key: "second",
					val: "b",
				},
				{
					act: actionDel,
					key: "first",
					val: "a",
				},
				{
					act: actionDel,
					key: "second",
					val: "b",
				},
			},
			expect: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := newContainer(false)
			for _, order := range test.do {
				if order.act == actionAdd {
					c.OnAdd(internal.KV{
						Key: order.key,
						Val: order.val,
					})
				} else {
					c.OnDelete(internal.KV{
						Key: order.key,
						Val: order.val,
					})
				}
			}
			assert.True(t, c.dirty.True())
			assert.ElementsMatch(t, test.expect, c.getValues())
			assert.False(t, c.dirty.True())
			assert.ElementsMatch(t, test.expect, c.getValues())
		})
	}
}
