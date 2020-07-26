package internal

import (
	"errors"
	"sort"
	"strconv"
	"testing"

	"zero/core/mathx"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestConsistent_addConn(t *testing.T) {
	b := NewConsistentBalancer(func(server string) (interface{}, error) {
		return mockConn{
			server: server,
		}, nil
	}, func(server string, conn interface{}) error {
		return errors.New("error")
	}, func(kv KV) string {
		return kv.Key
	})
	assert.Nil(t, b.AddConn(KV{
		Key: "thekey1",
		Val: "thevalue",
	}))
	assert.EqualValues(t, map[string]interface{}{
		"thekey1": mockConn{server: "thevalue"},
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
	assert.EqualValues(t, map[string]interface{}{
		"thekey1": mockConn{server: "thevalue"},
		"thekey2": mockConn{server: "thevalue"},
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
	assert.EqualValues(t, map[string]interface{}{
		"thekey2": mockConn{server: "thevalue"},
	}, b.conns)
	assert.EqualValues(t, map[string][]string{
		"thevalue": {"thekey2"},
	}, b.servers)
	assert.EqualValues(t, map[string]string{
		"thekey2": "thevalue",
	}, b.mapping)
	assert.False(t, b.IsEmpty())

	b.RemoveKey("thekey2")
	assert.Equal(t, 0, len(b.conns))
	assert.EqualValues(t, map[string][]string{}, b.servers)
	assert.EqualValues(t, map[string]string{}, b.mapping)
	assert.True(t, b.IsEmpty())
}

func TestConsistent_addConnError(t *testing.T) {
	b := NewConsistentBalancer(func(server string) (interface{}, error) {
		return nil, errors.New("error")
	}, func(server string, conn interface{}) error {
		return nil
	}, func(kv KV) string {
		return kv.Key
	})
	assert.NotNil(t, b.AddConn(KV{
		Key: "thekey1",
		Val: "thevalue",
	}))
	assert.Equal(t, 0, len(b.conns))
	assert.EqualValues(t, map[string][]string{}, b.servers)
	assert.EqualValues(t, map[string]string{}, b.mapping)
}

func TestConsistent_next(t *testing.T) {
	b := NewConsistentBalancer(func(server string) (interface{}, error) {
		return mockConn{
			server: server,
		}, nil
	}, func(server string, conn interface{}) error {
		return errors.New("error")
	}, func(kv KV) string {
		return kv.Key
	})
	b.initialize()

	_, ok := b.Next("any")
	assert.False(t, ok)

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
		val, ok := b.Next(strconv.Itoa(i))
		assert.True(t, ok)
		m[val]++
	}

	entropy := mathx.CalcEntropy(m, total)
	assert.Equal(t, size, len(m))
	assert.True(t, entropy > .95)

	for i := 0; i < size; i++ {
		b.RemoveKey("thekey/" + strconv.Itoa(i))
	}
	_, ok = b.Next()
	assert.False(t, ok)
}

func TestConsistentBalancer_Listener(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	b := NewConsistentBalancer(func(server string) (interface{}, error) {
		return mockConn{
			server: server,
		}, nil
	}, func(server string, conn interface{}) error {
		return nil
	}, func(kv KV) string {
		return kv.Key
	})
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
		sort.Strings(keys.([]string))
		sort.Strings(vals.([]string))
		assert.EqualValues(t, []string{"key1", "key2"}, keys)
		assert.EqualValues(t, []string{"val1", "val2"}, vals)
	})
	b.setListener(listener)
	b.notify("key2")
}

func TestConsistentBalancer_remove(t *testing.T) {
	b := NewConsistentBalancer(func(server string) (interface{}, error) {
		return mockConn{
			server: server,
		}, nil
	}, func(server string, conn interface{}) error {
		return nil
	}, func(kv KV) string {
		return kv.Key
	})

	assert.Nil(t, b.handlePrevious(nil))
	assert.Nil(t, b.handlePrevious([]string{"any"}))
}
