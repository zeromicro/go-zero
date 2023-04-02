package bloom

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis/redistest"
)

func TestRedisBitSet_New_Set_Test(t *testing.T) {
	store := redistest.CreateRedis(t)
	ctx := context.Background()

	bitSet := newRedisBitSet(store, "test_key", 1024)
	isSetBefore, err := bitSet.check(ctx, []uint{0})
	if err != nil {
		t.Fatal(err)
	}
	if isSetBefore {
		t.Fatal("Bit should not be set")
	}
	err = bitSet.set(ctx, []uint{512})
	if err != nil {
		t.Fatal(err)
	}
	isSetAfter, err := bitSet.check(ctx, []uint{512})
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
	store := redistest.CreateRedis(t)

	filter := New(store, "test_key", 64)
	assert.Nil(t, filter.Add([]byte("hello")))
	assert.Nil(t, filter.Add([]byte("world")))
	ok, err := filter.Exists([]byte("hello"))
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestFilter_Exists(t *testing.T) {
	store, clean := redistest.CreateRedisWithClean(t)

	rbs := New(store, "test", 64)
	_, err := rbs.Exists([]byte{0, 1, 2})
	assert.NoError(t, err)

	clean()
	rbs = New(store, "test", 64)
	_, err = rbs.Exists([]byte{0, 1, 2})
	assert.Error(t, err)
}

func TestRedisBitSet_check(t *testing.T) {
	store, clean := redistest.CreateRedisWithClean(t)
	ctx := context.Background()

	rbs := newRedisBitSet(store, "test", 0)
	assert.Error(t, rbs.set(ctx, []uint{0, 1, 2}))
	_, err := rbs.check(ctx, []uint{0, 1, 2})
	assert.Error(t, err)

	rbs = newRedisBitSet(store, "test", 64)
	_, err = rbs.check(ctx, []uint{0, 1, 2})
	assert.NoError(t, err)

	clean()
	rbs = newRedisBitSet(store, "test", 64)
	_, err = rbs.check(ctx, []uint{0, 1, 2})
	assert.Error(t, err)
}

func TestRedisBitSet_set(t *testing.T) {
	logx.Disable()
	store, clean := redistest.CreateRedisWithClean(t)
	ctx := context.Background()

	rbs := newRedisBitSet(store, "test", 0)
	assert.Error(t, rbs.set(ctx, []uint{0, 1, 2}))

	rbs = newRedisBitSet(store, "test", 64)
	assert.NoError(t, rbs.set(ctx, []uint{0, 1, 2}))

	clean()
	rbs = newRedisBitSet(store, "test", 64)
	assert.Error(t, rbs.set(ctx, []uint{0, 1, 2}))
}
