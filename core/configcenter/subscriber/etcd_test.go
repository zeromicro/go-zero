package subscriber

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/discov"
)

const (
	actionAdd = iota
	actionDel
)

func TestConfigCenterContainer(t *testing.T) {
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
			expect: []string(nil),
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
			expect: []string(nil),
		},
		{
			name: "add two, dup values",
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
			},
			expect: []string{"a"},
		},
		{
			name: "add three, dup values, delete two, add one",
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
					act: actionAdd,
					key: "forth",
					val: "c",
				},
			},
			expect: []string{"c"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var changed bool
			c := newContainer()
			c.AddListener(func() {
				changed = true
			})
			assert.Nil(t, c.GetValues())
			assert.False(t, changed)

			for _, order := range test.do {
				if order.act == actionAdd {
					c.OnAdd(discov.KV{
						Key: order.key,
						Val: order.val,
					})
				} else {
					c.OnDelete(discov.KV{
						Key: order.key,
						Val: order.val,
					})
				}
			}

			assert.True(t, changed)
			assert.ElementsMatch(t, test.expect, c.GetValues())
		})
	}
}
