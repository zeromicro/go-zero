package bloom

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stores/redis"
)

func TestRedisBitSet_New_Set_Test(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Error("Miniredis could not start")
	}
	defer s.Close()

	store := redis.NewRedis(s.Addr(), redis.NodeType)
	bitSet := newRedisBitSet(store, "test_key", 1024)
	isSetBefore, err := bitSet.check([]uint{0})
	if err != nil {
		t.Fatal(err)
	}
	if isSetBefore {
		t.Fatal("Bit should not be set")
	}
	err = bitSet.set([]uint{512})
	if err != nil {
		t.Fatal(err)
	}
	isSetAfter, err := bitSet.check([]uint{512})
	if err != nil {
		t.Fatal(err)
	}
	if !isSetAfter {
		t.Fatal("Bit should be set")
	}
	err = bitSet.expire(3600)
	if err != nil {
		t.Fatal(err)
	}
	err = bitSet.del()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRedisBitSet_Add(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Error("Miniredis could not start")
	}
	defer s.Close()

	store := redis.NewRedis(s.Addr(), redis.NodeType)
	filter := New(store, "test_key", 64)
	assert.Nil(t, filter.Add([]byte("hello")))
	assert.Nil(t, filter.Add([]byte("world")))
	ok, err := filter.Exists([]byte("hello"))
	assert.Nil(t, err)
	assert.True(t, ok)
}
