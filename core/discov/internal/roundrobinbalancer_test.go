package internal

import (
	"errors"
	"sort"
	"strconv"
	"testing"

	"zero/core/logx"
	"zero/core/mathx"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func init() {
	logx.Disable()
}

func TestRoundRobin_addConn(t *testing.T) {
	b := NewRoundRobinBalancer(func(server string) (interface{}, error) {
		return mockConn{
			server: server,
		}, nil
	}, func(server string, conn interface{}) error {
		return errors.New("error")
	}, false)
	assert.Nil(t, b.AddConn(KV{
		Key: "thekey1",
		Val: "thevalue",
	}))
	assert.EqualValues(t, []serverConn{
		{
			key:  "thekey1",
			conn: mockConn{server: "thevalue"},
		},
	}, b.conns)
	assert.EqualValues(t, map[string][]string{
		"thevalue": {"thekey1"},
	}, b.servers)
	assert.EqualValues(t, map[string]string{
		"thekey1": "thevalue",
	}, b.mapping)
	assert.Nil(t, b.AddConn(KV{
		Key: "thekey2",
		Val: "thevalue",
	}))
	assert.EqualValues(t, []serverConn{
		{
			key:  "thekey1",
			conn: mockConn{server: "thevalue"},
		},
		{
			key:  "thekey2",
			conn: mockConn{server: "thevalue"},
		},
	}, b.conns)
	assert.EqualValues(t, map[string][]string{
		"thevalue": {"thekey1", "thekey2"},
	}, b.servers)
	assert.EqualValues(t, map[string]string{
		"thekey1": "thevalue",
		"thekey2": "thevalue",
	}, b.mapping)
	assert.False(t, b.IsEmpty())

	b.RemoveKey("thekey1")
	assert.EqualValues(t, []serverConn{
		{
			key:  "thekey2",
			conn: mockConn{server: "thevalue"},
		},
	}, b.conns)
	assert.EqualValues(t, map[string][]string{
		"thevalue": {"thekey2"},
	}, b.servers)
	assert.EqualValues(t, map[string]string{
		"thekey2": "thevalue",
	}, b.mapping)
	assert.False(t, b.IsEmpty())

	b.RemoveKey("thekey2")
	assert.EqualValues(t, []serverConn{}, b.conns)
	assert.EqualValues(t, map[string][]string{}, b.servers)
	assert.EqualValues(t, map[string]string{}, b.mapping)
	assert.True(t, b.IsEmpty())
}

func TestRoundRobin_addConnExclusive(t *testing.T) {
	b := NewRoundRobinBalancer(func(server string) (interface{}, error) {
		return mockConn{
			server: server,
		}, nil
	}, func(server string, conn interface{}) error {
		return nil
	}, true)
	assert.Nil(t, b.AddConn(KV{
		Key: "thekey1",
		Val: "thevalue",
	}))
	assert.EqualValues(t, []serverConn{
		{
			key:  "thekey1",
			conn: mockConn{server: "thevalue"},
		},
	}, b.conns)
	assert.EqualValues(t, map[string][]string{
		"thevalue": {"thekey1"},
	}, b.servers)
	assert.EqualValues(t, map[string]string{
		"thekey1": "thevalue",
	}, b.mapping)
	assert.Nil(t, b.AddConn(KV{
		Key: "thekey2",
		Val: "thevalue",
	}))
	assert.EqualValues(t, []serverConn{
		{
			key:  "thekey2",
			conn: mockConn{server: "thevalue"},
		},
	}, b.conns)
	assert.EqualValues(t, map[string][]string{
		"thevalue": {"thekey2"},
	}, b.servers)
	assert.EqualValues(t, map[string]string{
		"thekey2": "thevalue",
	}, b.mapping)
	assert.False(t, b.IsEmpty())

	b.RemoveKey("thekey1")
	b.RemoveKey("thekey2")
	assert.EqualValues(t, []serverConn{}, b.conns)
	assert.EqualValues(t, map[string][]string{}, b.servers)
	assert.EqualValues(t, map[string]string{}, b.mapping)
	assert.True(t, b.IsEmpty())
}

func TestRoundRobin_addConnDupExclusive(t *testing.T) {
	b := NewRoundRobinBalancer(func(server string) (interface{}, error) {
		return mockConn{
			server: server,
		}, nil
	}, func(server string, conn interface{}) error {
		return errors.New("error")
	}, true)
	assert.Nil(t, b.AddConn(KV{
		Key: "thekey1",
		Val: "thevalue",
	}))
	assert.EqualValues(t, []serverConn{
		{
			key:  "thekey1",
			conn: mockConn{server: "thevalue"},
		},
	}, b.conns)
	assert.EqualValues(t, map[string][]string{
		"thevalue": {"thekey1"},
	}, b.servers)
	assert.EqualValues(t, map[string]string{
		"thekey1": "thevalue",
	}, b.mapping)
	assert.Nil(t, b.AddConn(KV{
		Key: "thekey",
		Val: "anothervalue",
	}))
	assert.Nil(t, b.AddConn(KV{
		Key: "thekey1",
		Val: "thevalue",
	}))
	assert.EqualValues(t, []serverConn{
		{
			key:  "thekey",
			conn: mockConn{server: "anothervalue"},
		},
		{
			key:  "thekey1",
			conn: mockConn{server: "thevalue"},
		},
	}, b.conns)
	assert.EqualValues(t, map[string][]string{
		"anothervalue": {"thekey"},
		"thevalue":     {"thekey1"},
	}, b.servers)
	assert.EqualValues(t, map[string]string{
		"thekey":  "anothervalue",
		"thekey1": "thevalue",
	}, b.mapping)
	assert.False(t, b.IsEmpty())

	b.RemoveKey("thekey")
	b.RemoveKey("thekey1")
	assert.EqualValues(t, []serverConn{}, b.conns)
	assert.EqualValues(t, map[string][]string{}, b.servers)
	assert.EqualValues(t, map[string]string{}, b.mapping)
	assert.True(t, b.IsEmpty())
}

func TestRoundRobin_addConnError(t *testing.T) {
	b := NewRoundRobinBalancer(func(server string) (interface{}, error) {
		return nil, errors.New("error")
	}, func(server string, conn interface{}) error {
		return nil
	}, true)
	assert.NotNil(t, b.AddConn(KV{
		Key: "thekey1",
		Val: "thevalue",
	}))
	assert.Nil(t, b.conns)
	assert.EqualValues(t, map[string][]string{}, b.servers)
	assert.EqualValues(t, map[string]string{}, b.mapping)
}

func TestRoundRobin_initialize(t *testing.T) {
	b := NewRoundRobinBalancer(func(server string) (interface{}, error) {
		return mockConn{
			server: server,
		}, nil
	}, func(server string, conn interface{}) error {
		return nil
	}, true)
	for i := 0; i < 100; i++ {
		assert.Nil(t, b.AddConn(KV{
			Key: "thekey/" + strconv.Itoa(i),
			Val: "thevalue/" + strconv.Itoa(i),
		}))
	}

	m := make(map[int]int)
	const total = 1000
	for i := 0; i < total; i++ {
		b.initialize()
		m[b.index]++
	}

	mi := make(map[interface{}]int, len(m))
	for k, v := range m {
		mi[k] = v
	}
	entropy := mathx.CalcEntropy(mi)
	assert.True(t, entropy > .95)
}

func TestRoundRobin_next(t *testing.T) {
	b := NewRoundRobinBalancer(func(server string) (interface{}, error) {
		return mockConn{
			server: server,
		}, nil
	}, func(server string, conn interface{}) error {
		return errors.New("error")
	}, true)
	const size = 100
	for i := 0; i < size; i++ {
		assert.Nil(t, b.AddConn(KV{
			Key: "thekey/" + strconv.Itoa(i),
			Val: "thevalue/" + strconv.Itoa(i),
		}))
	}

	m := make(map[interface{}]int)
	const total = 10000
	for i := 0; i < total; i++ {
		val, ok := b.Next()
		assert.True(t, ok)
		m[val]++
	}

	entropy := mathx.CalcEntropy(m)
	assert.Equal(t, size, len(m))
	assert.True(t, entropy > .95)

	for i := 0; i < size; i++ {
		b.RemoveKey("thekey/" + strconv.Itoa(i))
	}
	_, ok := b.Next()
	assert.False(t, ok)
}

func TestRoundRobinBalancer_Listener(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	b := NewRoundRobinBalancer(func(server string) (interface{}, error) {
		return mockConn{
			server: server,
		}, nil
	}, func(server string, conn interface{}) error {
		return nil
	}, true)
	assert.Nil(t, b.AddConn(KV{
		Key: "key1",
		Val: "val1",
	}))
	assert.Nil(t, b.AddConn(KV{
		Key: "key2",
		Val: "val2",
	}))

	listener := NewMockListener(ctrl)
	listener.EXPECT().OnUpdate(gomock.Any(), gomock.Any(), "key2").Do(func(keys, vals, _ interface{}) {
		sort.Strings(vals.([]string))
		sort.Strings(keys.([]string))
		assert.EqualValues(t, []string{"key1", "key2"}, keys)
		assert.EqualValues(t, []string{"val1", "val2"}, vals)
	})
	b.setListener(listener)
	b.notify("key2")
}

func TestRoundRobinBalancer_remove(t *testing.T) {
	b := NewRoundRobinBalancer(func(server string) (interface{}, error) {
		return mockConn{
			server: server,
		}, nil
	}, func(server string, conn interface{}) error {
		return nil
	}, true)

	assert.Nil(t, b.handlePrevious(nil, "any"))
	_, ok := b.doRemoveKv("any")
	assert.True(t, ok)
}
