package discov

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/discov/internal"
	"github.com/zeromicro/go-zero/core/stringx"
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
				c.AddListener(func() {
					changed = true
				})
				assert.Nil(t, c.GetValues())
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
				assert.ElementsMatch(t, test.expect, c.GetValues())
				assert.False(t, c.dirty.True())
				assert.ElementsMatch(t, test.expect, c.GetValues())
			})
		}
	}
}

func TestContainer_DuplicateAdd(t *testing.T) {
	c := newContainer(false)
	// Simulate 100 duplicate PUT events for the same key+value.
	for i := 0; i < 100; i++ {
		c.OnAdd(internal.KV{Key: "etcd-key", Val: "host:1234"})
	}
	assert.ElementsMatch(t, []string{"host:1234"}, c.GetValues())
	// Internal slice must not have grown beyond one entry.
	c.lock.Lock()
	assert.Len(t, c.values["host:1234"], 1)
	c.lock.Unlock()
}

func TestContainer_KeyValueChange(t *testing.T) {
	c := newContainer(false)
	c.OnAdd(internal.KV{Key: "etcd-key", Val: "host:1234"})
	assert.ElementsMatch(t, []string{"host:1234"}, c.GetValues())

	// Key moves to a different server value.
	c.OnAdd(internal.KV{Key: "etcd-key", Val: "host:5678"})
	assert.ElementsMatch(t, []string{"host:5678"}, c.GetValues())

	// Old server must be fully removed; a subsequent delete must leave nothing.
	c.OnDelete(internal.KV{Key: "etcd-key", Val: "host:5678"})
	assert.Empty(t, c.GetValues())
}

// TestContainer_ExclusiveMode verifies that adding successive keys for the same
// value in exclusive mode leaves only the latest key and evicts all prior ones.
func TestContainer_ExclusiveMode(t *testing.T) {
	c := newContainer(true)
	c.OnAdd(internal.KV{Key: "key1", Val: "server1"})
	c.OnAdd(internal.KV{Key: "key2", Val: "server1"})
	c.OnAdd(internal.KV{Key: "key3", Val: "server1"})

	assert.ElementsMatch(t, []string{"server1"}, c.GetValues())
	c.lock.Lock()
	assert.Equal(t, []string{"key3"}, c.values["server1"], "only the latest key must remain")
	assert.NotContains(t, c.mapping, "key1", "key1 must have been evicted")
	assert.NotContains(t, c.mapping, "key2", "key2 must have been evicted")
	assert.Equal(t, "server1", c.mapping["key3"])
	c.lock.Unlock()
}

// TestContainer_ExclusiveMode_MultipleEvictions injects 3 keys for the same
// value directly into internal state and then triggers the exclusive eviction
// loop via OnAdd. This exercises the range-over-previous fix: iterating over
// the live slice (range keys) would corrupt iteration when doRemoveKey
// compacts the shared underlying array in-place, causing the second and third
// keys to be skipped; ranging over the deep copy (range previous) is safe.
func TestContainer_ExclusiveMode_MultipleEvictions(t *testing.T) {
	c := newContainer(true)

	// Bypass the exclusive invariant to simulate 3 pre-existing keys for the
	// same value — the state that would expose the in-place aliasing bug.
	c.lock.Lock()
	c.values["server1"] = []string{"key1", "key2", "key3"}
	c.mapping["key1"] = "server1"
	c.mapping["key2"] = "server1"
	c.mapping["key3"] = "server1"
	c.lock.Unlock()

	// Adding key4 must evict all three existing keys via the exclusive loop.
	c.OnAdd(internal.KV{Key: "key4", Val: "server1"})

	assert.ElementsMatch(t, []string{"server1"}, c.GetValues())
	c.lock.Lock()
	assert.Equal(t, []string{"key4"}, c.values["server1"], "all prior keys must be evicted")
	assert.NotContains(t, c.mapping, "key1", "key1 must be evicted")
	assert.NotContains(t, c.mapping, "key2", "key2 must be evicted")
	assert.NotContains(t, c.mapping, "key3", "key3 must be evicted")
	assert.Equal(t, "server1", c.mapping["key4"])
	c.lock.Unlock()
}

func TestSubscriber(t *testing.T) {
	sub := new(Subscriber)
	Exclusive()(sub)
	c := newContainer(sub.exclusive)
	WithContainer(c)(sub)
	sub.items = c
	var count int32
	sub.AddListener(func() {
		atomic.AddInt32(&count, 1)
	})
	c.notifyChange()
	assert.Empty(t, sub.Values())
	assert.Equal(t, int32(1), atomic.LoadInt32(&count))
}

func TestWithSubEtcdAccount(t *testing.T) {
	endpoints := []string{"localhost:2379"}
	user := stringx.Rand()
	WithSubEtcdAccount(user, "bar")(&Subscriber{
		endpoints: endpoints,
	})
	account, ok := internal.GetAccount(endpoints)
	assert.True(t, ok)
	assert.Equal(t, user, account.User)
	assert.Equal(t, "bar", account.Pass)
}

func TestWithExactMatch(t *testing.T) {
	sub := new(Subscriber)
	WithExactMatch()(sub)
	c := newContainer(sub.exclusive)
	sub.items = c
	var count int32
	sub.AddListener(func() {
		atomic.AddInt32(&count, 1)
	})
	c.notifyChange()
	assert.Empty(t, sub.Values())
	assert.Equal(t, int32(1), atomic.LoadInt32(&count))
}

func TestSubscriberClose(t *testing.T) {
	l := newContainer(false)
	sub := &Subscriber{
		endpoints: []string{"localhost:12379"},
		key:       "foo",
		items:     l,
	}
	assert.NotPanics(t, func() {
		sub.Close()
	})
}
