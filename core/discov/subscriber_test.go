package discov

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/discov/internal"
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
				},
			},
			expect: []string{"b"},
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
				},
				{
					act: actionDel,
					key: "second",
				},
			},
			expect: []string{},
		},
		{
			name: "add three, dup values, delete two",
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
					act: actionAdd,
					key: "third",
					val: "a",
				},
				{
					act: actionDel,
					key: "first",
				},
				{
					act: actionDel,
					key: "second",
				},
			},
			expect: []string{"a"},
		},
		{
			name: "add three, dup values, delete two, delete not added",
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
					act: actionAdd,
					key: "third",
					val: "a",
				},
				{
					act: actionDel,
					key: "first",
				},
				{
					act: actionDel,
					key: "second",
				},
				{
					act: actionDel,
					key: "forth",
				},
			},
			expect: []string{"a"},
		},
	}

	exclusives := []bool{true, false}
	for _, test := range tests {
		for _, exclusive := range exclusives {
			t.Run(test.name, func(t *testing.T) {
				var changed bool
				c := newContainer(exclusive)
				c.addListener(func() {
					changed = true
				})
				assert.Nil(t, c.getValues())
				assert.False(t, changed)

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

				assert.True(t, changed)
				assert.True(t, c.dirty.True())
				assert.ElementsMatch(t, test.expect, c.getValues())
				assert.False(t, c.dirty.True())
				assert.ElementsMatch(t, test.expect, c.getValues())
			})
		}
	}
}
